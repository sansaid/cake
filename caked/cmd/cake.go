package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"sync"
	"time"

	"github.com/sansaid/cake/caked/pb"
	"github.com/sansaid/cake/utils"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	dockerClient "github.com/docker/docker/client"
	specs "github.com/opencontainers/image-spec/specs-go/v1"
	log "github.com/sirupsen/logrus"
)

type SliceCtx struct {
	Container *pb.Slice
	WaitGroup *sync.WaitGroup
	Cancel    context.CancelFunc
}

type Slices struct {
	sync.RWMutex
	slices map[string]SliceCtx
}

type Cake struct {
	pb.UnimplementedCakedServer
	DockerClient  ModContainerAPIClient
	HttpClient    ModHTTPClient
	RunningSlices Slices
	StopTimeout   time.Duration
}

type CakeOpts func(*Cake)

// has type CakeOpts
func WithStopTimeout(duration time.Duration) CakeOpts {
	return func(c *Cake) {
		c.StopTimeout = duration
	}
}

func NewCake(opts ...CakeOpts) *Cake {
	client := utils.Must(dockerClient.NewClientWithOpts()).(ModContainerAPIClient)

	cake := &Cake{
		DockerClient: client,
		HttpClient:   &http.Client{Timeout: time.Second * 3},
		StopTimeout:  30 * time.Second,
	}

	for _, opt := range opts {
		opt(cake)
	}

	return cake
}

// gRPC server methods - this should only get called by the gRPC client (should never be called directly in this code)
func (c *Cake) RunSlice(ctx context.Context, container *pb.Slice) (*pb.SliceStatus, error) {
	log.Info("Starting slice for image %s with tag %s", container.ImageName, container.Tag)

	var wg *sync.WaitGroup
	ctx, cancel := context.WithCancel(ctx)

	wg.Add(1)
	go c.PollSlice(ctx, wg, container, 5*time.Second)

	c.RunningSlices.Lock()
	c.RunningSlices.slices[container.ImageName] = SliceCtx{
		Container: container,
		WaitGroup: wg,
		Cancel:    cancel,
	}
	c.RunningSlices.Unlock()

	return &pb.SliceStatus{
		Status:  0,
		Message: "Container successfully started",
	}, nil
}

// gRPC server methods - this should only get called by the gRPC client (should never be called directly in this code)
func (c *Cake) StopSlice(ctx context.Context, container *pb.Slice) (*pb.SliceStatus, error) {
	fmt.Printf("Stopping slice: %#v", container) // TODO - NEXT: implement StopSlice gRPC operation

	return &pb.SliceStatus{
		Status:  0,
		Message: "Container successfully stopped",
	}, nil
}

func (c *Cake) PollSlice(ctx context.Context, wg *sync.WaitGroup, slice *pb.Slice, frequency time.Duration) {
	for {
		select {
		case <-ctx.Done():
			// Only terminating latest digest since this is assumed to be called synchronously after any previous digests
			// have been terminated as a result of an update
			c.TermDigest(slice.ImageName, slice.LatestDigest)
			wg.Done()
			return
		default:
			if c.UpdateLatestDigest(slice) {
				if err := c.RunDigest(slice.ImageName, slice.LatestDigest, slice); err != nil {
					log.Errorf("Could not run latest digest %s for image %s: %w", slice.LatestDigest, slice.ImageName, err)
				}

				if err := c.TermDigest(slice.ImageName, slice.PreviousDigest); err != nil {
					log.Errorf("Could not terminate previous digest %s for image %s: %w", slice.PreviousDigest, slice.ImageName, err)
				}
			}

			time.Sleep(frequency)
		}
	}
}

func (c *Cake) UpdateLatestDigest(slice *pb.Slice) bool {
	slice.LastChecked = time.Now().Unix()

	latestDigest, latestDigestTime, err := c.GetLatestDigest(slice)

	if err != nil {
		log.Errorf("Could not update latest digest for image %s: %w", slice.ImageName, err)
		return false
	}

	if slice.LatestDigest != latestDigest {
		slice.PreviousDigest = slice.LatestDigest
		slice.PreviousDigestTime = slice.LatestDigestTime
		slice.LatestDigest = latestDigest
		slice.LatestDigestTime = latestDigestTime

		return true
	}

	return false
}

func (c *Cake) RunDigest(imageName string, digest string, slice *pb.Slice) error {
	if digest != "" {
		ctx := context.TODO()
		hostConfig := container.HostConfig{}
		networkConfig := network.NetworkingConfig{}

		platformSpecs := &specs.Platform{
			Architecture: slice.Architecture,
			OS:           slice.OS,
		}

		containerConfig := container.Config{
			Image: fmt.Sprintf("%s@%s", slice.ImageName, slice.LatestDigest),
		}

		createdContainer, err := c.DockerClient.ContainerCreate(ctx, &containerConfig, &hostConfig, &networkConfig, platformSpecs, "")

		if err != nil {
			return fmt.Errorf("Could not create container for image %s and digest %s: %w", imageName, digest, err)
		}

		err = c.DockerClient.ContainerStart(ctx, createdContainer.ID, types.ContainerStartOptions{})

		if err != nil {
			return fmt.Errorf("Could not start container for image %s and digest %s: %w", imageName, digest, err)
		}

		// TODO: do we want to wait for container to be started? Will only be for logging purposes - debugging will have to be done through Docker
	}

	log.Infof("No digest to run for image %s", imageName)
	return nil
}

func (c *Cake) TermDigest(imageName string, digest string) error {
	if digest != "" {
		containerIds, err := c.ListRunningContainerIds(imageName, digest) // TODO: do you want termination to be on all running containers running this digest? Shouldn't you separate what's managed by Cake and what was manually spun up?

		if err != nil {
			return fmt.Errorf("Could not stop container. Error in listing containers: %w", err)
		}

		for _, id := range containerIds {
			ctx := context.TODO()
			err := c.DockerClient.ContainerStop(ctx, id, &c.StopTimeout)

			if err != nil {
				log.Errorf("Could not issue stop to container %s: %w", id, err)
				continue
			}

			// TODO: see if you want to wait for the stop command to finish - there isn't really any further action to take other than to log the stopping failed
		}
	}

	log.Infof("No digest to terminate for image %s", imageName)
	return nil
}

func (c *Cake) ListRunningContainerIds(image string, digest string) ([]string, error) {
	containerList, err := c.DockerClient.ContainerList(context.TODO(), types.ContainerListOptions{All: false})

	if err != nil {
		return []string{}, err
	}

	imageName := fmt.Sprintf("%s@%s", image, digest)
	containerIds := []string{}

	for _, containerInstance := range containerList {
		if containerInstance.Image == imageName {
			containerIds = append(containerIds, containerInstance.ID)
		}
	}

	return containerIds, nil
}

func (c *Cake) GetLatestDigest(slice *pb.Slice) (string, int64, error) {
	repoURL := fmt.Sprintf("https://registry.hub.docker.com/v2/repositories/%s/tags?ordering=last_updated", slice.ImageName)

	repoList := RepoList{}

	if err := c.MarshalHttp(repoURL, &repoList); err != nil {
		return "", 0, err
	}

	imagesForArch := []Image{}
	images := repoList.Results[0].Images

	for _, image := range images {
		if image.Architecture == string(slice.Architecture) {
			imagesForArch = append(imagesForArch, image)
		}
	}

	sort.Sort(ByLastPushedDesc(imagesForArch))

	latestImage := imagesForArch[0]
	latestImageDigest := latestImage.Digest
	latestImagePushTime := latestImage.LastPushed

	latestDigestTime := latestImagePushTime.Unix()

	return latestImageDigest, latestDigestTime, nil
}

func (c *Cake) MarshalHttp(url string, t interface{}) error {
	// Cannot test unhappy path, only happy path
	req, err := http.NewRequest(http.MethodGet, url, nil)

	if err != nil {
		return fmt.Errorf("Could not perform get request on %s: %w", url, err)
	}

	resp, err := c.HttpClient.Do(req)
	defer resp.Body.Close()

	if err != nil {
		return fmt.Errorf("Could not read request response from %s: %w", url, err)
	}

	// Could not test unhappy path, only happy path
	err = json.NewDecoder(resp.Body).Decode(t)

	if err != nil {
		return fmt.Errorf("Could not decode JSON from URL %s: %w", url, err) // TODO: decide if you want
	}

	return nil
}
