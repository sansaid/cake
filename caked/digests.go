package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	dockerClient "github.com/docker/docker/client"
	specs "github.com/opencontainers/image-spec/specs-go/v1"
	pb "github.com/sansaid/cake/pb"
	"github.com/sansaid/cake/utils"
)

// Private functions are declared with vars to allow for stubbing in tests - most of these functions are not unit testable because they're side effect heavy
var getImageDigests = func(repoURL string, arch string) []Image {
	// TODO: see if filter for architecture can be done through REST API
	repoList := RepoList{}
	decodeResponse(repoURL, &repoList)

	archImages := []Image{}
	images := repoList.Results[0].Images

	for _, image := range images {
		if image.Architecture == string(arch) {
			archImages = append(archImages, image)
		}
	}

	return archImages
}

var getRunningContainerIds = func(client *dockerClient.Client, image string, digest string) []string {
	containerList := listRunningContainers(client)

	imageName := fmt.Sprintf("%s@%s", image, digest)
	containerIds := []string{}

	for _, containerInstance := range containerList {
		if containerInstance.Image == imageName {
			containerIds = append(containerIds, containerInstance.ID)
		}
	}

	return containerIds
}

var pullImage = func(client *dockerClient.Client, imageRef string) {
	reader, err := client.ImagePull(context.TODO(), imageRef, types.ImagePullOptions{})
	defer reader.Close()

	utils.Check(err, "Could not pull image")
}

var listImages = func(client *dockerClient.Client) []types.ImageSummary {
	imageList, err := client.ImageList(context.TODO(), types.ImageListOptions{})

	if err != nil {
		panic(err)
	}

	return imageList
}

var listRunningContainers = func(client *dockerClient.Client) []types.Container {
	containerList, err := client.ContainerList(context.TODO(), types.ContainerListOptions{All: false})

	utils.Check(err, "Could not list running containers")

	return containerList
}

var createContainer = func(client *dockerClient.Client, cakeContainer *pb.Container, containerConfig container.Config, hostConfig container.HostConfig, networkConfig network.NetworkingConfig) string {
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

var decodeResponse = func(url string, t interface{}) interface{} {
	resp, err := http.Get(url)
	defer resp.Body.Close()

	utils.Check(err, fmt.Sprintf("Could not perform get request on %s", url))

	err = json.NewDecoder(resp.Body).Decode(t)

	utils.Check(err, fmt.Sprintf("Could not decode JSON", url))

	return t
}

// GetLatestDigest - get the latest digest for the container specified
func (c *Cake) GetLatestDigest(cakeContainer *pb.Container) *Cake {
	cakeContainer.LastChecked = time.Now().Unix()

	repoURL := fmt.Sprintf("https://registry.hub.docker.com/v2/repositories/%s/tags?ordering=last_updated", cakeContainer.ImageName)
	images := getImageDigests(repoURL, cakeContainer.Architecture)

	sort.Sort(ByLastPushedDesc(images))

	imageLatestDigest := images[0].Digest
	imageLastPushedTime := images[0].LastPushed

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

func (c *Cake) StopPreviousDigest(cakeContainer *pb.Container) {
	if cakeContainer.PreviousDigest != "" {
		containerIds := getRunningContainerIds(c.DockerClient, cakeContainer.ImageName, cakeContainer.PreviousDigest)

		for _, id := range containerIds {
			stopContainer(c, id)

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

		pullImage(c.DockerClient, imageRef)
	}
}

func (c *Cake) IsLatestDigestPulled(cakeContainer *pb.Container) bool {
	imageList := listImages(c.DockerClient)

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

		id := createContainer(c.DockerClient, cakeContainer, containerConfig, hostConfig, networkConfig)

		c.ContainersRunning.Lock()
		c.ContainersRunning.containers[id] = 0
		c.ContainersRunning.Unlock()
	} else {
		runningContainers := getRunningContainerIds(c.DockerClient, cakeContainer.ImageName, cakeContainer.LatestDigest)

		// Checks if the container is in Cake's control. If it's not,
		// Cake adds it to its list of running containers.
		for _, id := range runningContainers {
			c.ContainersRunning.RLock()
			_, managedByCake := c.ContainersRunning.containers[id]
			c.ContainersRunning.RLock()

			if !(managedByCake) {
				c.ContainersRunning.Lock()
				c.ContainersRunning.containers[id] = 0
				c.ContainersRunning.Unlock()
			}
		}
	}
}

func (c *Cake) IsLatestDigestRunning(cakeContainer *pb.Container) bool {
	containerList := listRunningContainers(c.DockerClient)

	latestImageName := fmt.Sprintf("%s@%s", cakeContainer.ImageName, cakeContainer.LatestDigest)

	for _, containerInstance := range containerList {
		if latestImageName == containerInstance.Image {
			return true
		}
	}

	return false
}
