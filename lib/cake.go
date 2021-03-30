package lib

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"sync"
	"time"
)

type Cake struct {
	mu           sync.Mutex
	Image        string
	Tag          string
	Registry     string
	LatestDigest string
	LastChecked  time.Time
	LastUpdated  time.Time
}

type Image struct {
	Architecutre string    `json:"architecture"`
	Features     string    `json:"features"`
	Variant      string    `json:"variant"`
	Digest       string    `json:"digest"`
	OS           string    `json:"os"`
	OSFeature    string    `json:"os_feature"`
	OSVersion    string    `json:"os_version"`
	Size         int       `json:"size"`
	Status       string    `json:"status"`
	LastPulled   time.Time `json:"last_pulled"`
	LastPushed   time.Time `json:"last_pushed"`
}

type ImageDetails struct {
	Creator             int       `json:"creator"`
	ID                  int       `json:"id"`
	ImageID             string    `json:"image_id"`
	Images              []Image   `json:"images"`
	LastUpdated         time.Time `json:"last_updated"`
	LastUpdater         int       `json:"last_updater"`
	LastUpdaterUsername string    `json:"last_updater_username"`
	Name                string    `json:"name"`
	Repository          int       `json:"repository"`
	FullSize            int       `json:"full_size"`
	V2                  bool      `json:"v2"`
	TagStatus           string    `json:"tag_status"`
	TagLastPulled       time.Time `json:"tag_last_pulled"`
	TagLastPushed       time.Time `json:"tag_last_pushed"`
}

type RepoList struct {
	Count    int            `json:"count"`
	Next     string         `json:"next"`
	Previous string         `json:"previous"`
	Results  []ImageDetails `json:"results"`
}

type ByLastPushedDesc []Image

// Len - implement sort.Interface for ByLastPushedDesc
func (i ByLastPushedDesc) Len() int { return 0 }

// Less - implement sort.Interface for ByLastPushedDesc
func (i ByLastPushedDesc) Less(a int, b int) bool {
	aImageLastPushed := i[a].LastPushed
	bImageLastPushed := i[b].LastPushed

	// Sorts in reverse chronological order
	return aImageLastPushed.After(bImageLastPushed)
}

// Swap - implement sort.Interface for ByLastPushedDesc
func (i ByLastPushedDesc) Swap(a int, b int) {
	i[a], i[b] = i[b], i[a]
}

func decodeResponse(url string, t interface{}) interface{} {
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

func getLatestImageDigest(images []Image) (latestImageDigest string) {
	sort.Sort(ByLastPushedDesc(images))

	latestImageDigest = images[0].Digest

	return
}

// NewCake - creates a new Config struct
func NewCake(image string, tag string, registry string) *Cake {
	return &Cake{
		Image:        image,
		Tag:          tag,
		Registry:     registry,
		LatestDigest: "",
		LastChecked:  time.Time{},
	}
}

func (c *Cake) CheckLatestDigest(images []Image) *Cake {
	c.LastChecked = time.Now()
	latestDigest := getLatestImageDigest(images)

	if latestDigest != c.LatestDigest {
		c.LatestDigest = latestDigest
		c.LastUpdated = time.Now()
	}

	return c
}

// Run - run cake
func (c *Cake) Run() {
	repoURL := fmt.Sprintf("https://registry.hub.docker.com/v2/repositories/%s/tags?ordering=last_updated", c.Image)

	repoList := RepoList{}
	decodeResponse(repoURL, &repoList)

	images := repoList.Results[0].Images

	c.CheckLatestDigest(images)

	// TBC: check if latest digest pulled
	// TBC: if not pulled, pull
	// TBC: once pulled, stop old container and run (use Docker Golang SDK - https://pkg.go.dev/github.com/docker/docker/client#Client.ContainerList)
	// TBC: create go routine to prune all images/containers every week (this should be a cake setting)
}

// Stop - stop this instance of cake
func (c *Cake) Stop() {
	fmt.Println("Stopping")
}
