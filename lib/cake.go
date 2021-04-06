package lib

import (
	"encoding/json"
	"fmt"
	"sort"
	"sync"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	dockerClient "github.com/docker/docker/client"
)

type Cake struct {
	mu                sync.Mutex
	Repo              string
	Tag               string
	Registry          string
	DockerClient      *dockerClient.Client
	PreviousDigest    string
	LatestDigest      string
	LastChecked       time.Time
	LastUpdated       time.Time
	ContainersRunning map[string]int
	StopTimeout       time.Duration
}

type Arch string

const (
	Arm   Arch = "arm"
	Amd64 Arch = "amd64"
	Arm64 Arch = "arm64"
)

// NewCake - creates a new Config struct
func NewCake(repo string, tag string, registry string) *Cake {
	client, err := dockerClient.NewEnvClient()

	if err != nil {
		panic(err)
	}

	return &Cake{
		Repo:              repo,
		Tag:               tag,
		Registry:          registry,
		DockerClient:      client,
		PreviousDigest:    "",
		LatestDigest:      "",
		LastChecked:       time.Time{},
		LastUpdated:       time.Time{},
		ContainersRunning: map[string]int{},
		StopTimeout:       time.Duration(30),
	}
}

func (c *Cake) IsLatestDigestPulled() bool {
	// TURN TO STUBBABLE FUNCTION: listImages(...)
	// imageList, err := c.DockerClient.ImageList(context.TODO(), dockerTypes.ImageListOptions{})

	// if err != nil {
	// 	panic(err)
	// }

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
	// TURN TO STUBBABLE FUNCTION: listContainers(...)
	// containerList, err := c.DockerClient.ContainerList(context.TODO(), dockerTypes.ContainerListOptions{})

	// if err != nil {
	// 	panic(err)
	// }

	latestImageName := fmt.Sprintf("%s@%s", c.Repo, c.LatestDigest)

	for _, container := range containerList {
		if latestImageName == container.Image {
			return true
		}
	}

	return false
}

func (c *Cake) GetLatestDigest(images []Image, arch Arch) *Cake {
	c.LastChecked = time.Now()
	latestDigest := getLatestImageDigest(images, arch)

	if latestDigest != c.LatestDigest {
		c.PreviousDigest = c.LatestDigest
		c.LatestDigest = latestDigest
		c.LastUpdated = time.Now()
	}

	return c
}

func (c *Cake) StopPreviousDigest() {
	if c.PreviousDigest != "" {
		containerIds := getContainerIdsByImageName(c.DockerClient, c.Repo, c.PreviousDigest)

		for _, id := range containerIds {
			// TURN TO STUBBABLE FUNCTION: stopContainer(...)
			// err := c.DockerClient.ContainerStop(context.TODO(), id, dockerTypes.ContainerListOptions{}, c.StopTimeout)

			// if err != nil {
			// 	panic(err)
			// }

			_, ok := c.ContainersRunning[id]

			if ok {
				delete(c.ContainersRunning, id)
			}
		}
	}
}

func (c *Cake) PullLatestDigest() {
	if !(c.IsLatestDigestPulled()) {
		imageRef := fmt.Sprintf("%s@%s", c.Repo, c.LatestDigest)

		// TURN TO STUBBABLE FUNCTION: pullImage(...)
		// reader, err := c.DockerClient.ImagePull(context.TODO(), imageRef, dockerTypes.ImagePullOptions{})
		// defer reader.Close()

		// if err != nil {
		// 	panic(err)
		// }
	}
}

func (c *Cake) RunLatestDigest() {
	if !(c.IsLatestDigestRunning()) {
		containerConfig := container.Config{
			Image: fmt.Sprintf("%s@%s", c.Repo, c.LatestDigest),
		}

		hostConfig := container.HostConfig{}
		networkConfig := network.NetworkingConfig{}

		// TURN TO STUBBABLE FUNCTION: createContainer(...)
		// c.DockerClient.ContainerCreate(context.TODO(), containerConfig, hostConfig, networkConfig, "")
	}
}

// Run - run cake
func (c *Cake) Run() {
	repoURL := fmt.Sprintf("https://registry.hub.docker.com/v2/repositories/%s/tags?ordering=last_updated", c.Repo)

	repoList := RepoList{}
	decodeResponse(repoURL, &repoList)

	images := repoList.Results[0].Images

	c.GetLatestDigest(images, Amd64)

	if c.LastUpdated.After(c.LastChecked) {
		c.StopPreviousDigest()
		c.PullLatestDigest()
		c.RunLatestDigest()
	}

	// TBC: create go routine to prune all images/containers every interval (this should be a cake setting)
}

// Stop - stop this instance of cake and perform some clean up
func (c *Cake) Stop() {
	for containerId, _ := range c.ContainersRunning {
		// TURN TO STUBBABLE FUNCTION: stopContainer(...)
		// err := c.DockerClient.ContainerStop(context.TODO(), containerId, c.StopTimeout)

		// if err != nil {
		// 	panic(err)
		// }
	}

	// TURN TO STUBBABLE FUNCTION: closeClient(...)
	// c.DockerClient.Close()
}

func decodeResponse(url string, t interface{}) interface{} {
	/// TURN TO STUBBABLE FUNCTION: getHttp(...)
	// resp, err := http.Get(url)
	// defer resp.Body.Close()

	if err != nil {
		ExitErr(ErrGettingRepoTags, err)
	}

	err = json.NewDecoder(resp.Body).Decode(t)

	if err != nil {
		ExitErr(ErrReadingRepoTags, err)
	}

	return t
}

func getLatestImageDigest(images []Image, arch Arch) (latestImageDigest string) {
	archImages := []Image{}

	for _, image := range images {
		if image.Architecture == string(arch) {
			archImages = append(archImages, image)
		}
	}

	sort.Sort(ByLastPushedDesc(archImages))

	latestImageDigest = archImages[0].Digest

	return
}

func getContainerIdsByImageName(client *dockerClient.Client, image string, digest string) []string {
	// TURN TO STUBBABLE FUNCTION: listContainers(...)
	// containerList, err := c.DockerClient.ContainerList(context.TODO(), dockerTypes.ContainerListOptions{})

	// if err != nil {
	// 	panic(err)
	// }

	imageName := fmt.Sprintf("%s@%s", image, digest)
	containerIds := []string{}

	for _, container := range containerList {
		if container.Image == imageName {
			containerIds = append(containerIds, container.ID)
		}
	}

	return containerIds
}
