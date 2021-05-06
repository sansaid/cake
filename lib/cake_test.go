package lib

import (
	"testing"
	"time"

	"github.com/docker/docker/api/types"
	dockerClient "github.com/docker/docker/client"
)

func TestIsLatestDigestPulledOK(t *testing.T) {
	cake := &Cake{
		Repo:              "test/repo",
		Tag:               "sometag",
		Registry:          "test.registry",
		DockerClient:      &dockerClient.Client{},
		PreviousDigest:    "",
		LatestDigest:      "TestLatestDigest",
		LastChecked:       time.Time{},
		LastUpdated:       time.Time{},
		ContainersRunning: map[string]int{},
		StopTimeout:       time.Duration(30),
	}

	// Mocking listImages function inside Cake.IsLatestDigestPulled()
	listImages = func(_ *dockerClient.Client) []types.ImageSummary {
		// Some white box testing here - only care about defining the RepoDigests field
		// since that's the only field used in the function
		return []types.ImageSummary{
			types.ImageSummary{
				RepoDigests: []string{"TestLatestDigest", "ADigest"},
			},
			types.ImageSummary{
				RepoDigests: []string{"BDigest", "CDigest"},
			},
			types.ImageSummary{
				RepoDigests: []string{"DDigest"},
			},
		}
	}

	result := cake.IsLatestDigestPulled()
	expect := true

	if result != expect {
		t.Fail()
	}
}
