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

const (
	SLICE_STATUS_OK = iota
	SLICE_STATUS_NOT_FOUND
)

type SliceCtx struct {
	Container *pb.Slice // TODO: change name to Slice - Containers will be reserved for list of container IDs the slice manages (we may want to separate the containers managed by cake and those managed externally)
	WaitGroup *sync.WaitGroup
	Cancel    context.CancelFunc
}

type Slices struct {
	sync.RWMutex
	slices map[string]SliceCtx
}

type Cake struct {
	pb.UnimplementedCakedServer
	DockerClient  CakeContainerAPIClient
	HttpClient    CakeHTTPClient
	RunningSlices *Slices
	StopTimeout   time.Duration
}

type CakeOpts func(*Cake)

func WithStopTimeout(duration time.Duration) CakeOpts {
	return func(c *Cake) {
		c.StopTimeout = duration
	}
}

func WithDockerClient(client CakeContainerAPIClient) CakeOpts {
	return func(c *Cake) {
		c.DockerClient = client
	}
}

func WithHttpClient(client CakeHTTPClient) CakeOpts {
	return func(c *Cake) {
		c.HttpClient = client
	}
}

func NewCake(opts ...CakeOpts) *Cake {
	client := utils.Must(dockerClient.NewClientWithOpts()).(CakeContainerAPIClient)

	cake := &Cake{
		DockerClient:  client,
		HttpClient:    &http.Client{Timeout: time.Second * 3},
		StopTimeout:   30 * time.Second,
		RunningSlices: &Slices{slices: make(map[string]SliceCtx)},
	}

	for _, opt := range opts {
		opt(cake)
	}

	return cake
}

func (c *Cake) AddNewSlice(slice *pb.Slice, wg *sync.WaitGroup, cancelFunc context.CancelFunc) {
	c.RunningSlices.RLock()
	c.RunningSlices.slices[slice.ImageName] = SliceCtx{
		Container: slice,
		WaitGroup: wg,
		Cancel:    cancelFunc,
	}
	c.RunningSlices.RUnlock()
}

// gRPC server methods - this should only get called by the gRPC client (should never be called directly in this code)
func (c *Cake) RunSlice(ctx context.Context, slice *pb.Slice) (*pb.SliceStatus, error) {
	log.Infof("Starting slice for image %s", slice.ImageName)

	// TODO: create a c.ValidateSlice(slice) function

	var wg *sync.WaitGroup
	ctx, cancel := context.WithCancel(ctx)

	c.AddNewSlice(slice, wg, cancel)

	wg.Add(1)
	go c.PollSlice(ctx, wg, slice, 5*time.Second)

	return &pb.SliceStatus{
		Status:  SLICE_STATUS_OK,
		Message: "Container successfully started", // TODO: modify status message to include container IDs
	}, nil
}

// gRPC server methods - this should only get called by the gRPC client (should never be called directly in this code)
func (c *Cake) StopSlice(ctx context.Context, image *pb.Image) (*pb.SliceStatus, error) {
	imageName := image.Name

	log.Infof("Stopping slice with image name %s", imageName)

	if _, ok := c.RunningSlices.slices[imageName]; !ok {
		log.Errorf("StopSlice command issued for non-existent slice: %s", imageName)

		return &pb.SliceStatus{
			Status:  SLICE_STATUS_NOT_FOUND,
			Message: fmt.Sprintf("Slice for image %s cannot be found", imageName), // TODO: modify status message to include container IDs
		}, nil
	}

	c.RunningSlices.slices[imageName].Cancel()         // TODO: figure out how to deal with timeouts here - can these methods be transactional?
	c.RunningSlices.slices[imageName].WaitGroup.Wait() // TODO: figure out how to deal with timeouts here - can these methods be transactional?

	return &pb.SliceStatus{
		Status:  0,
		Message: "Container successfully stopped", // TODO: modify status message to include container IDs
	}, nil
}

func (c *Cake) PollSlice(ctx context.Context, wg *sync.WaitGroup, slice *pb.Slice, frequency time.Duration) {
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			// Only terminating latest digest since this is assumed to be called synchronously after any previous digests
			// have been terminated as a result of an update
			c.TermDigest(slice.ImageName, slice.LatestDigest)

			// Delegating responsibility of removing the slice reference here so that calling the cancel method externally automatically removes
			// the reference in Cake
			c.RunningSlices.RLock()
			delete(c.RunningSlices.slices, slice.ImageName)
			c.RunningSlices.RUnlock()

			return
		default:
			if c.UpdateLatestDigest(slice) {
				if err := c.RunDigest(slice.ImageName, slice.LatestDigest, slice); err != nil {
					log.Errorf("Could not run latest digest %s for image %s: %s", slice.LatestDigest, slice.ImageName, err)
				}

				if err := c.TermDigest(slice.ImageName, slice.PreviousDigest); err != nil {
					log.Errorf("Could not terminate previous digest %s for image %s: %s", slice.PreviousDigest, slice.ImageName, err)
				}
			}

			time.Sleep(frequency) // TODO: polling frequency should be configured by slice
		}
	}
}

// Keeping this as a Cake method since it relies on the Docker client wrapped inside Cake
func (c *Cake) UpdateLatestDigest(slice *pb.Slice) bool {
	slice.LastChecked = time.Now().Unix()

	latestDigest, latestDigestTime, err := c.GetLatestDigest(slice)

	if err != nil {
		log.Errorf("Could not update latest digest for image %s: %s", slice.ImageName, err)
		return false
	}

	if latestDigest == "" && latestDigestTime == 0 {
		log.Debug("No latest image matching platform requirements (architecture and OS)")
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

// TODO: make this a Slice method
func (c *Cake) RunDigest(imageName string, digest string, slice *pb.Slice) error {
	if digest != "" {
		ctx := context.TODO()
		hostConfig := container.HostConfig{}
		networkConfig := network.NetworkingConfig{}

		platformSpecs := &specs.Platform{
			Architecture: slice.Architecture,
			OS:           slice.Os,
		}

		containerConfig := container.Config{
			Image: fmt.Sprintf("%s@%s", slice.ImageName, slice.LatestDigest),
		}

		createdContainer, err := c.DockerClient.ContainerCreate(ctx, &containerConfig, &hostConfig, &networkConfig, platformSpecs, "")

		if err != nil {
			return fmt.Errorf("could not create container for image %s and digest %s: %w", imageName, digest, err)
		}

		err = c.DockerClient.ContainerStart(ctx, createdContainer.ID, types.ContainerStartOptions{})

		if err != nil {
			return fmt.Errorf("could not start container for image %s and digest %s: %w", imageName, digest, err)
		}

		// TODO: do we want to wait for container to be started? Will only be for logging purposes - debugging will have to be done through Docker
	}

	log.Infof("No digest to run for image %s", imageName)
	return nil
}

// TODO: make this a Slice method
func (c *Cake) TermDigest(imageName string, digest string) error {
	if digest != "" {
		containerIds, err := c.ListRunningContainerIds(imageName, digest) // TODO: do you want termination to be on all running containers running this digest? Shouldn't you separate what's managed by Cake and what was manually spun up?

		if err != nil {
			return fmt.Errorf("could not stop container. Error in listing containers: %w", err)
		}

		for _, id := range containerIds {
			ctx := context.TODO()
			err := c.DockerClient.ContainerStop(ctx, id, &c.StopTimeout)

			if err != nil {
				log.Errorf("Could not issue stop to container %s: %s", id, err)
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

// TODO: make this a Slice method
func (c *Cake) GetLatestDigest(slice *pb.Slice) (string, int64, error) {
	repoURL := fmt.Sprintf("https://registry.hub.docker.com/v2/repositories/%s/tags?ordering=last_updated", slice.ImageName) // TODO: container image should be fully qualified - error if not (this will allow pulling from arbitrary repos)

	repoList := RepoList{}

	if err := c.MarshalHttp(repoURL, &repoList); err != nil {
		return "", 0, err
	}

	compatibleImages := []Image{}

	if len(repoList.Results) == 0 {
		return "", 0, fmt.Errorf("response returned Results with length 0 - could not pull latest image")
	}

	allImages := repoList.Results[0].Images // assumes there will always be at least one image in the repo

	if len(allImages) == 0 {
		return "", 0, fmt.Errorf("response returned latest image with Images list that has length 0 - could not pull latest image")
	}

	for _, image := range allImages {
		if missingRequiredFields(image) {
			log.Errorf("One of the images returned does not contain the required fields - something is wrong with the response. Bad image: %#v.", image)
			continue
		}

		if image.Architecture == string(slice.Architecture) && image.OS == string(slice.Os) {
			compatibleImages = append(compatibleImages, image)
		}
	}

	if len(compatibleImages) == 0 {
		return "", 0, nil
	}

	sort.Sort(ByLastPushedDesc(compatibleImages))

	latestImage := compatibleImages[0]
	latestImageDigest := latestImage.Digest
	latestImagePushTime := latestImage.LastPushed

	latestDigestTime := latestImagePushTime.Unix()

	return latestImageDigest, latestDigestTime, nil
}

// Made into a method of Cake so that we can stub the method with a dummy implementation in testing
func (c *Cake) MarshalHttp(url string, t interface{}) error {
	// Cannot test unhappy path, only happy path
	req, err := http.NewRequest(http.MethodGet, url, nil)

	if err != nil {
		return fmt.Errorf("could not perform get request on %s: %w", url, err)
	}

	resp, err := c.HttpClient.Do(req)

	if err != nil {
		return fmt.Errorf("could not read request response from %s: %w", url, err)
	}

	defer resp.Body.Close()

	// Could not test unhappy path, only happy path
	err = json.NewDecoder(resp.Body).Decode(t)

	if err != nil {
		return fmt.Errorf("could not decode JSON from URL %s: %w", url, err)
	}

	return nil
}

func (c *Cake) Stop(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	for _, slice := range c.RunningSlices.slices {
		slice.Cancel()
	}

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("timed out stopping cake: %w", ctx.Err())
		default:
			if len(c.RunningSlices.slices) == 0 {
				return nil
			}
		}
	}
}

func missingRequiredFields(image Image) bool {
	return image.LastPushed.IsZero() || image.Digest == "" || image.OS == "" || image.Architecture == ""
}
