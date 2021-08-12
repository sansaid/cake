package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"sync"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	dockerClient "github.com/docker/docker/client"
	specs "github.com/opencontainers/image-spec/specs-go/v1"
	pb "github.com/sansaid/cake/pb"
	"github.com/sansaid/cake/utils"
)

type ContainerList struct {
	sync.RWMutex
	containers map[string]struct{}
}

type Cake struct {
	pb.UnimplementedCakedServer
	DockerClient      ModContainerAPIClient
	HttpClient        ModHTTPClient
	ContainersRunning ContainerList
	StopTimeout       time.Duration
}

func NewCake() *Cake {
	client, err := dockerClient.NewEnvClient()

	utils.Check(err, "Cannot create cake client")

	return &Cake{
		DockerClient: client,
		HttpClient:   &http.Client{Timeout: time.Second * 3},
		StopTimeout:  30 * time.Second,
	}
}

func (c *Cake) ListRunningContainerIds(image string, digest string) []string {
	containerList, err := c.DockerClient.ContainerList(context.TODO(), types.ContainerListOptions{All: false})

	utils.Check(err, "Could not list running containers")

	imageName := fmt.Sprintf("%s@%s", image, digest)
	containerIds := []string{}

	for _, containerInstance := range containerList {
		if containerInstance.Image == imageName {
			containerIds = append(containerIds, containerInstance.ID)
		}
	}

	return containerIds
}

func (c *Cake) CreateContainer(client ModContainerAPIClient, cakeContainer *pb.Container, containerConfig container.Config, hostConfig container.HostConfig, networkConfig network.NetworkingConfig) string {
	ctx := context.TODO()

	platform := &specs.Platform{
		Architecture: cakeContainer.Architecture,
		OS:           cakeContainer.OS,
	}

	createdContainer, err := client.ContainerCreate(ctx, &containerConfig, &hostConfig, &networkConfig, platform, "")

	utils.Check(err, "Could not create container")

	err = client.ContainerStart(ctx, createdContainer.ID, types.ContainerStartOptions{})

	utils.Check(err, "Could not start container")

	statC, errC := client.ContainerWait(ctx, createdContainer.ID, "created")

	select {
	case err := <-errC:
		utils.Check(err, "Error waiting for container to be created")
		return ""
	case stat := <-statC:
		if stat.StatusCode == 0 {
			return createdContainer.ID
		} else {
			utils.Check(err, fmt.Sprintf("Error waiting for container to be started - Docker exit code %d", stat.StatusCode))
			return ""
		}
	}
}

func (c *Cake) Get(url string, t interface{}) interface{} {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	defer req.Body.Close()

	utils.Check(err, fmt.Sprintf("Could not perform get request on %s", url))

	resp, err := c.HttpClient.Do(req)
	defer resp.Body.Close()

	utils.Check(err, fmt.Sprintf("Could not read request response from %s", url))

	err = json.NewDecoder(resp.Body).Decode(t)

	utils.Check(err, fmt.Sprintf("Could not decode JSON from URL %s", url))

	return t
}

// GetLatestDigest - get the latest digest for the container specified
func (c *Cake) GetLatestDigest(cakeContainer *pb.Container) *Cake {
	cakeContainer.LastChecked = time.Now().Unix()

	repoURL := fmt.Sprintf("https://registry.hub.docker.com/v2/repositories/%s/tags?ordering=last_updated", cakeContainer.ImageName)

	repoList := RepoList{}
	c.Get(repoURL, &repoList)

	archImages := []Image{}
	images := repoList.Results[0].Images

	for _, image := range images {
		if image.Architecture == string(cakeContainer.Architecture) {
			archImages = append(archImages, image)
		}
	}

	sort.Sort(ByLastPushedDesc(images))

	imageLatestDigest := archImages[0].Digest
	imageLastPushedTime := archImages[0].LastPushed

	containerLatestDigestTime := time.Unix(cakeContainer.LatestDigestTime, 0)

	// Is it certain that every new image push will create a new digest?
	// Can an existing digest be updated?
	if imageLastPushedTime.After(containerLatestDigestTime) {
		cakeContainer.PreviousDigest = cakeContainer.LatestDigest
		cakeContainer.PreviousDigestTime = cakeContainer.LatestDigestTime
		cakeContainer.LatestDigest = imageLatestDigest
		cakeContainer.LatestDigestTime = imageLastPushedTime.Unix()
		cakeContainer.LastUpdated = time.Now().Unix()
	}

	return c
}

func (c *Cake) stopContainer(id string) {
	ctx := context.TODO()

	err := c.DockerClient.ContainerStop(ctx, id, &c.StopTimeout)

	utils.Check(err, "Could not issue stop to container")

	_, errC := c.DockerClient.ContainerWait(ctx, id, "removed")

	select {
	case err := <-errC:
		utils.Check(err, "Error waiting for container to be removed")
	default:
		return
	}
}

func (c *Cake) StopPreviousDigest(cakeContainer *pb.Container) {
	if cakeContainer.PreviousDigest != "" {
		containerIds := c.ListRunningContainerIds(c.DockerClient, cakeContainer.ImageName, cakeContainer.PreviousDigest)

		for _, id := range containerIds {
			c.stopContainer(id)

			c.ContainersRunning.RLock()
			_, exists := c.ContainersRunning.containers[id]
			c.ContainersRunning.RUnlock()

			if exists {
				c.ContainersRunning.Lock()
				delete(c.ContainersRunning.containers, id)
				c.ContainersRunning.Unlock()
			}
		}
	}
}

func (c *Cake) PullLatestDigest(cakeContainer *pb.Container) {
	if !(c.IsLatestDigestPulled(cakeContainer)) {
		imageRef := fmt.Sprintf("%s@%s", cakeContainer.ImageName, cakeContainer.LatestDigest)

		reader, err := c.DockerClient.ImagePull(context.TODO(), imageRef, types.ImagePullOptions{})
		defer reader.Close()

		utils.Check(err, "Could not pull image")
	}
}

func (c *Cake) IsLatestDigestPulled(cakeContainer *pb.Container) bool {
	imageList, err := c.DockerClient.ImageList(context.TODO(), types.ImageListOptions{})

	if err != nil {
		panic(err)
	}

	latestImageName := fmt.Sprintf("%s@%s", cakeContainer.ImageName, cakeContainer.LatestDigest)

	var digestList []string

	// Extract digests
	for _, image := range imageList {
		digestList = append(digestList, image.RepoDigests...)
	}

	// Find matches to digest
	for _, digest := range digestList {
		if latestImageName == digest {
			return true
		}
	}

	return false
}

func (c *Cake) RunLatestDigest(cakeContainer *pb.Container) {
	if !(c.IsLatestDigestRunning(cakeContainer)) {
		containerConfig := container.Config{
			Image: fmt.Sprintf("%s@%s", cakeContainer.ImageName, cakeContainer.LatestDigest),
		}

		hostConfig := container.HostConfig{}
		networkConfig := network.NetworkingConfig{}

		id := c.CreateContainer(c.DockerClient, cakeContainer, containerConfig, hostConfig, networkConfig)

		c.ContainersRunning.Lock()
		c.ContainersRunning.containers[id] = struct{}{}
		c.ContainersRunning.Unlock()
	} else {
		runningContainers := c.ListRunningContainerIds(c.DockerClient, cakeContainer.ImageName, cakeContainer.LatestDigest)

		// Checks if the container is in Cake's control. If it's not,
		// Cake adds it to its list of running containers.
		for _, id := range runningContainers {
			c.ContainersRunning.RLock()
			_, managedByCake := c.ContainersRunning.containers[id]
			c.ContainersRunning.RLock()

			if !(managedByCake) {
				c.ContainersRunning.Lock()
				c.ContainersRunning.containers[id] = struct{}{}
				c.ContainersRunning.Unlock()
			}
		}
	}
}

func (c *Cake) IsLatestDigestRunning(cakeContainer *pb.Container) bool {
	containerList, err := c.DockerClient.ContainerList(context.TODO(), types.ContainerListOptions{All: false})

	utils.Check(err, "Could not list running containers")

	latestImageName := fmt.Sprintf("%s@%s", cakeContainer.ImageName, cakeContainer.LatestDigest)

	for _, containerInstance := range containerList {
		if latestImageName == containerInstance.Image {
			return true
		}
	}

	return false
}
