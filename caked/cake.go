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
)

// TODO: The use of the word `digest` in variable and function names is inconsistently being used - need to make this more consistent
// TODO: Write functionality to sync `cake` with the local system that's being managed externally to it
// TODO: Think about how to deal with pruning containers and images with `cake` on RasbPi - should stopped containers also be deleted and have their images removed?

type Arch string

const (
	Arm   Arch = "arm"
	Amd64 Arch = "amd64"
	Arm64 Arch = "arm64"
)

var listImages = func(client *dockerClient.Client) []types.ImageSummary {
	imageList, err := client.ImageList(context.TODO(), types.ImageListOptions{})

	if err != nil {
		panic(err)
	}

	return imageList
}

var listContainers = func(client *dockerClient.Client) []types.Container {
	containerList, err := client.ContainerList(context.TODO(), types.ContainerListOptions{})

	if err != nil {
		panic(err)
	}

	return containerList
}

var stopContainer = func(cake *Cake, id string) {
	ctx := context.TODO()

	err := cake.DockerClient.ContainerStop(ctx, id, &cake.StopTimeout)

	if err != nil {
		panic(err)
	}

	_, errC := cake.DockerClient.ContainerWait(ctx, id, "removed")

	select {
	case err := <-errC:
		panic(err)
	default:
		return
	}
}

var closeClient = func(c *Cake) {
	c.DockerClient.Close()
}

var decodeResponse = func(url string, t interface{}) interface{} {
	resp, err := http.Get(url)
	defer resp.Body.Close()

	if err != nil {
		ExitErr(ErrGettingRepoTags, err)
	}

	err = json.NewDecoder(resp.Body).Decode(t)

	if err != nil {
		ExitErr(ErrReadingRepoTags, err)
	}

	return t
}

var getImageDigests = func(repoURL string, arch Arch) []Image {
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
	containerList := listContainers(client)

	imageName := fmt.Sprintf("%s@%s", image, digest)
	containerIds := []string{}

	for _, container := range containerList {
		if container.Image == imageName {
			containerIds = append(containerIds, container.ID)
		}
	}

	return containerIds
}

var pullImage = func(client *dockerClient.Client, imageRef string) {
	reader, err := client.ImagePull(context.TODO(), imageRef, types.ImagePullOptions{})
	defer reader.Close()

	if err != nil {
		panic(err)
	}
}

var createContainer = func(client *dockerClient.Client, containerConfig container.Config, hostConfig container.HostConfig, networkConfig network.NetworkingConfig) string {
	ctx := context.TODO()

	platform := &specs.Platform{
		Architecture: "amd64",
		OS:           "linux",
	}

	createdContainer, err := client.ContainerCreate(ctx, &containerConfig, &hostConfig, &networkConfig, platform, "")

	if err != nil {
		panic(err)
	}

	err = client.ContainerStart(ctx, createdContainer.ID, types.ContainerStartOptions{})

	if err != nil {
		panic(err)
	}

	statC, errC := client.ContainerWait(ctx, createdContainer.ID, "created")

	select {
	case err := <-errC:
		panic(err)
	case stat := <-statC:
		if stat.StatusCode == 0 {
			return createdContainer.ID
		} else {
			ExitErr(ErrCreateContainer, fmt.Errorf("could not start container: exit code %d", stat.StatusCode))
			return ""
		}
	}
}

func (c *Cake) IsLatestDigestPulled() bool {
	imageList := listImages(c.DockerClient)

	latestImageName := fmt.Sprintf("%s@%s", c.Repo, c.LatestDigest)

	digestList := []string{}

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

func (c *Cake) IsLatestDigestRunning() bool {
	containerList := listContainers(c.DockerClient)

	latestImageName := fmt.Sprintf("%s@%s", c.Repo, c.LatestDigest)

	for _, container := range containerList {
		if latestImageName == container.Image {
			return true
		}
	}

	return false
}

func (c *Cake) GetLatestDigest(arch Arch) *Cake {
	c.LastChecked = time.Now()

	repoURL := fmt.Sprintf("https://registry.hub.docker.com/v2/repositories/%s/tags?ordering=last_updated", c.Repo)
	images := getImageDigests(repoURL, arch)

	sort.Sort(ByLastPushedDesc(images))

	latestDigest := images[0].Digest
	latestDigestTime := images[0].LastPushed

	// Is it certain that every new image push will create a new digest?
	// Can an existing digest be updated?
	if latestDigestTime.After(c.LatestDigestTime) {
		c.PreviousDigest = c.LatestDigest
		c.PreviousDigestTime = c.LatestDigestTime
		c.LatestDigest = latestDigest
		c.LatestDigestTime = latestDigestTime
		c.LastUpdated = time.Now()
	}

	return c
}

func (c *Cake) StopPreviousDigest() {
	if c.PreviousDigest != "" {
		containerIds := getRunningContainerIds(c.DockerClient, c.Repo, c.PreviousDigest)

		for _, id := range containerIds {
			stopContainer(c, id)

			_, wasRunning := c.ContainersRunning[id]

			if wasRunning {
				delete(c.ContainersRunning, id)
			}
		}
	}
}

func (c *Cake) PullLatestDigest() {
	if !(c.IsLatestDigestPulled()) {
		imageRef := fmt.Sprintf("%s@%s", c.Repo, c.LatestDigest)

		pullImage(c.DockerClient, imageRef)
	}
}

func (c *Cake) RunLatestDigest() {
	if !(c.IsLatestDigestRunning()) {
		containerConfig := container.Config{
			Image: fmt.Sprintf("%s@%s", c.Repo, c.LatestDigest),
		}

		hostConfig := container.HostConfig{}
		networkConfig := network.NetworkingConfig{}

		id := createContainer(c.DockerClient, containerConfig, hostConfig, networkConfig)

		c.ContainersRunning[id] = 0
	} else {
		runningContainers := getRunningContainerIds(c.DockerClient, c.Repo, c.LatestDigest)

		// Checks if the container is in Cake's control. If it's not,
		// Cake adds it to its list of running containers.
		for _, id := range runningContainers {
			_, managedByCake := c.ContainersRunning[id]

			if !(managedByCake) {
				c.ContainersRunning[id] = 0
			}
		}
	}
}

// Run - run cake
func (c *Cake) Run() {
	c.GetLatestDigest(Amd64)

	if c.LastUpdated.After(c.LastChecked) {
		c.StopPreviousDigest()
		c.PullLatestDigest()
		c.RunLatestDigest()
	}
}

// Stop - stop this instance of cake and perform some clean up
func (c *Cake) Stop() {
	for containerId, _ := range c.ContainersRunning {
		stopContainer(c, containerId)
	}

	defer closeClient(c)
}
