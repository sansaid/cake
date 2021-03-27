package lib

import (
	"encoding/json"
	"fmt"
	"json"
	"net/http"
	"sync"

	"github.com/sansaid/cake/cmds"
)

type Cake struct {
	mu    sync.Mutex
	image *string
}

type Image struct {
	architecutre string `json:"architecture"`
	features     string `json:"features"`
	variant      string `json:"variant"`
	digest       string `json:"digest"`
	os           string `json:"os"`
	os_feature   string `json:"os_feature"`
	os_version   string `json:"os_version"`
	size         int    `json:"size"`
	status       string `json:"status"`
	lastPulled   string `json:"last_pulled"`
	lastPushed   string `json:"last_pushed"`
}

type ImageDetails struct {
	creator             int     `json:"creator"`
	id                  int     `json:"id"`
	imageId             string  `json:"image_id"`
	images              []Image `json:"images"`
	lastUpdated         string  `json:"last_updated"`
	lastUpdater         int     `json:"last_updater"`
	lastUpdaterUsername string  `json:"last_updater_username"`
	name                string  `json:"name"`
	repository          int     `json:"repository"`
	fullSize            int     `json:"full_size"`
	v2                  bool    `json:"v2"`
	tagStatus           string  `json:"tag_status"`
	tagLastPulled       string  `json:"tag_last_pulled"`
	tagLastPushed       string  `json:"tag_last_pushed"`
}

func decodeResponse(url string, t interface{}) {
	resp, err := http.Get(url)
	defer resp.Body.Close()

	if err != nil {
		cmds.ExitErr(cmds.ErrGettingRepoTags)
	}

	if err != nil {
		cmds.ExitErr(cmds.ErrReadingRepoTags)
	}

	return json.NewDecoder(resp.Body).Decode(t)
}

// NewCake - creates a new Config struct
func NewCake(image *string) *Cake {
	return &Cake{
		image: image,
	}
}

// Run - run cake
func (c *Cake) Run() {
	repoURL := fmt.Sprintf("https://registry.hub.docker.com/v2/repositories/%s/tags?ordering=last_updated", c.image)

	var imageDetails *ImageDetails

	for {
		decodeResponse(repoURL, imageDetails)
	}
}
