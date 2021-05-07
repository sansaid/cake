package lib

import (
	"testing"
	"time"

	"github.com/docker/docker/api/types"
	dockerClient "github.com/docker/docker/client"
)

func TestIsLatestDigestPulled_OK(t *testing.T) {
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
			{
				RepoDigests: []string{"test/repo@TestLatestDigest", "ADigest"},
			},
			{
				RepoDigests: []string{"BDigest", "CDigest"},
			},
			{
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

func TestIsLatestDigestPulled_Bad(t *testing.T) {
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
			{
				RepoDigests: []string{"TestLatestDigest", "ADigest"},
			},
			{
				RepoDigests: []string{"BDigest", "CDigest"},
			},
			{
				RepoDigests: []string{"DDigest"},
			},
		}
	}

	result := cake.IsLatestDigestPulled()
	expect := false

	if result != expect {
		t.Fail()
	}
}

func TestIsLatestDigestRunning_OK(t *testing.T) {
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

	// Mocking listContainers function inside Cake.IsLatestDigestRunning()
	listContainers = func(_ *dockerClient.Client) []types.Container {
		// Some white box testing here - only care about defining the Image field
		// since that's the only field used in the function
		return []types.Container{
			{
				Image: "test/repo@TestLatestDigest",
			},
			{
				Image: "BDigest",
			},
			{
				Image: "DDigest",
			},
		}
	}

	result := cake.IsLatestDigestRunning()
	expect := true

	if result != expect {
		t.Fail()
	}
}

func TestIsLatestDigestRunning_Bad(t *testing.T) {
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

	// Mocking listContainer function inside Cake.IsLatestDigestPulled()
	listContainers = func(_ *dockerClient.Client) []types.Container {
		// Some white box testing here - only care about defining the Image field
		// since that's the only field used in the function
		return []types.Container{
			{
				Image: "TestLatestDigest",
			},
			{
				Image: "BDigest",
			},
			{
				Image: "DDigest",
			},
		}
	}

	result := cake.IsLatestDigestRunning()
	expect := false

	if result != expect {
		t.Fail()
	}
}
