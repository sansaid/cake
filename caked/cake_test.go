package main

import (
	"testing"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	dockerClient "github.com/docker/docker/client"
)

type ResultMap struct {
	result string
	expect string
}

func log(t *testing.T, result interface{}, expect interface{}) {
	t.Logf("Expected %v, got %v", expect, result)
}

// Test when latest digest is pulled
func TestIsLatestDigestPulled_OK(t *testing.T) {
	cake := &Cake{
		Repo:           "test/repo",
		PreviousDigest: "",
		LatestDigest:   "TestLatestDigest",
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
		log(t, result, expect)
		t.Fail()
	}
}

// Test when latest digest is not pulled
func TestIsLatestDigestPulled_Bad(t *testing.T) {
	cake := &Cake{
		Repo:           "test/repo",
		PreviousDigest: "",
		LatestDigest:   "TestLatestDigest",
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
		log(t, result, expect)
		t.Fail()
	}
}

// Test when latest digest is running
func TestIsLatestDigestRunning_OK(t *testing.T) {
	cake := &Cake{
		Repo:           "test/repo",
		PreviousDigest: "",
		LatestDigest:   "TestLatestDigest",
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
		log(t, result, expect)
		t.Fail()
	}
}

// Test when latest digest is not running
func TestIsLatestDigestRunning_Bad(t *testing.T) {
	cake := &Cake{
		Repo:           "test/repo",
		PreviousDigest: "",
		LatestDigest:   "TestLatestDigest",
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
		log(t, result, expect)
		t.Fail()
	}
}

// Test when current latest is not in the list of returned digests
func TestGetLatestDigest_LatestNotListed(t *testing.T) {
	cake := &Cake{
		Repo:               "test/repo",
		LatestDigest:       "TestLatestDigest",
		LatestDigestTime:   time.Date(2020, time.May, 10, 0, 0, 0, 0, *&time.UTC),
		PreviousDigest:     "TestPreviousDigest",
		PreviousDigestTime: time.Date(2020, time.May, 9, 0, 0, 0, 0, *&time.UTC),
	}

	curr_time := time.Now()
	old_last_updated := cake.LastUpdated

	getImageDigests = func(_ string, _ Arch) []Image {
		return []Image{
			{
				Digest:     "OldDigestOne",
				LastPushed: time.Date(2020, time.May, 7, 0, 0, 0, 0, *&time.UTC),
			},
			{
				Digest:     "OldDigestTwo",
				LastPushed: time.Date(2020, time.May, 6, 0, 0, 0, 0, *&time.UTC),
			},
			{
				Digest:     "OldDigestThree",
				LastPushed: time.Date(2020, time.May, 5, 0, 0, 0, 0, *&time.UTC),
			},
		}
	}

	cake.GetLatestDigest(Amd64)

	results := []ResultMap{
		{
			result: cake.LatestDigest,
			expect: "TestLatestDigest",
		},
		{
			result: cake.PreviousDigest,
			expect: "TestPreviousDigest",
		},
	}

	for _, r := range results {
		if r.result != r.expect {
			log(t, r.result, r.expect)
			t.Fail()
		}
	}

	if cake.LastChecked.Before(curr_time) {
		t.Logf("Expected cake.LastChecked to be after current time, got: %v", cake.LastChecked)
		t.Fail()
	}

	if cake.LastUpdated != old_last_updated {
		log(t, old_last_updated, cake.LastUpdated)
		t.Fail()
	}
}

// Test when current latest is the latest in returned digests
func TestGetLatestDigest_CurrLatestIsLatest(t *testing.T) {
	cake := &Cake{
		Repo:               "test/repo",
		LatestDigest:       "TestLatestDigest",
		LatestDigestTime:   time.Date(2020, time.May, 10, 0, 0, 0, 0, *&time.UTC),
		PreviousDigest:     "TestPreviousDigest",
		PreviousDigestTime: time.Date(2020, time.May, 9, 0, 0, 0, 0, *&time.UTC),
	}

	curr_time := time.Now()
	old_last_updated := cake.LastUpdated

	getImageDigests = func(_ string, _ Arch) []Image {
		return []Image{
			{
				Digest:     "TestPreviousDigest",
				LastPushed: time.Date(2020, time.May, 9, 0, 0, 0, 0, *&time.UTC),
			},
			{
				Digest:     "TestLatestDigest",
				LastPushed: time.Date(2020, time.May, 10, 0, 0, 0, 0, *&time.UTC),
			},
			{
				Digest:     "OldDigestThree",
				LastPushed: time.Date(2020, time.May, 5, 0, 0, 0, 0, *&time.UTC),
			},
		}
	}

	cake.GetLatestDigest(Amd64)

	results := []ResultMap{
		{
			result: cake.LatestDigest,
			expect: "TestLatestDigest",
		},
		{
			result: cake.PreviousDigest,
			expect: "TestPreviousDigest",
		},
	}

	for _, r := range results {
		if r.result != r.expect {
			log(t, r.result, r.expect)
			t.Fail()
		}
	}

	if cake.LastChecked.Before(curr_time) {
		t.Logf("Expected cake.LastChecked to be after current time, got: %v", cake.LastChecked)
		t.Fail()
	}

	if cake.LastUpdated != old_last_updated {
		log(t, old_last_updated, cake.LastUpdated)
		t.Fail()
	}
}

// Test when latest is not the one in cake
func TestGetLatestDigest_CurrLatestIsNotLatest(t *testing.T) {
	cake := &Cake{
		Repo:               "test/repo",
		LatestDigest:       "TestCurrentDigest",
		LatestDigestTime:   time.Date(2020, time.May, 10, 0, 0, 0, 0, *&time.UTC),
		PreviousDigest:     "TestPreviousDigest",
		PreviousDigestTime: time.Date(2020, time.May, 9, 0, 0, 0, 0, *&time.UTC),
	}

	curr_time := time.Now()

	getImageDigests = func(_ string, _ Arch) []Image {
		return []Image{
			{
				Digest:     "TestPreviousDigest",
				LastPushed: time.Date(2020, time.May, 9, 0, 0, 0, 0, *&time.UTC),
			},
			{
				Digest:     "TestCurrentDigest",
				LastPushed: time.Date(2020, time.May, 10, 0, 0, 0, 0, *&time.UTC),
			},
			{
				Digest:     "TestLatestDigest",
				LastPushed: time.Date(2020, time.May, 11, 0, 0, 0, 0, *&time.UTC),
			},
		}
	}

	cake.GetLatestDigest(Amd64)

	results := []ResultMap{
		{
			result: cake.LatestDigest,
			expect: "TestLatestDigest",
		},
		{
			result: cake.LatestDigestTime.String(),
			expect: time.Date(2020, time.May, 11, 0, 0, 0, 0, *&time.UTC).String(),
		},
		{
			result: cake.PreviousDigest,
			expect: "TestCurrentDigest",
		},
		{
			result: cake.PreviousDigestTime.String(),
			expect: time.Date(2020, time.May, 10, 0, 0, 0, 0, *&time.UTC).String(),
		},
	}

	for _, r := range results {
		if r.result != r.expect {
			log(t, r.result, r.expect)
			t.Fail()
		}
	}

	if cake.LastChecked.Before(curr_time) {
		t.Logf("Expected cake.LastChecked to be after current time, got: %v", cake.LastChecked)
		t.Fail()
	}

	if cake.LastUpdated.Before(curr_time) {
		t.Logf("Expected cake.LastUpdated to be after current time, got: %v", cake.LastUpdated)
		t.Fail()
	}
}

// Test stopping when no previous digest specified
func TestStopPreviousDigest_NoPreviouDigest(t *testing.T) {
	var containerStopped string

	getRunningContainerIds = func(_ *dockerClient.Client, _ string, _ string) []string {
		return []string{}
	}

	stopContainer = func(_ *Cake, id string) {
		containerStopped = id
	}

	cake := Cake{
		ContainersRunning: map[string]int{
			"a": 0,
			"b": 0,
			"c": 0,
		},
	}

	expected := map[string]int{
		"a": 0,
		"b": 0,
		"c": 0,
	}

	cake.StopPreviousDigest()

	for id := range expected {
		if _, ok := cake.ContainersRunning[id]; ok != true {
			t.Logf("Expected %s to be running, but instead was stopped", id)
			t.Fail()
		}
	}

	if containerStopped != "" {
		log(t, containerStopped, "")
		t.Fail()
	}
}

// Test stopping when previous digest is not running
func TestStopPreviousDigest_PreviouDigestNotRunning(t *testing.T) {
	var containerStopped string

	getRunningContainerIds = func(_ *dockerClient.Client, _ string, _ string) []string {
		return []string{}
	}

	stopContainer = func(_ *Cake, id string) {
		containerStopped = id
	}

	cake := Cake{
		PreviousDigest: "TestPreviousDigest",
		ContainersRunning: map[string]int{
			"a": 0,
			"b": 0,
			"c": 0,
		},
	}

	expected := map[string]int{
		"a": 0,
		"b": 0,
		"c": 0,
	}

	cake.StopPreviousDigest()

	for id := range expected {
		if _, ok := cake.ContainersRunning[id]; ok != true {
			t.Logf("Expected %s to be running, but instead was stopped", id)
			t.Fail()
		}
	}

	if containerStopped != "" {
		log(t, containerStopped, "")
		t.Fail()
	}
}

// Test stopping when previous digest is running
func TestStopPreviousDigest_PreviouDigestIsRunning(t *testing.T) {
	containerStopped := []string{}

	getRunningContainerIds = func(_ *dockerClient.Client, _ string, _ string) []string {
		return []string{"a", "c"}
	}

	stopContainer = func(_ *Cake, id string) {
		containerStopped = append(containerStopped, id)
	}

	cake := Cake{
		PreviousDigest: "TestPreviousDigest",
		ContainersRunning: map[string]int{
			"a": 0,
			"b": 0,
			"c": 0,
			"d": 0,
		},
	}

	expected := map[string]int{
		"b": 0,
		"d": 0,
	}

	cake.StopPreviousDigest()

	for id := range expected {
		if _, ok := cake.ContainersRunning[id]; ok != true {
			t.Logf("Expected %s to be running, but instead was stopped", id)
			t.Fail()
		}
	}

	if len(containerStopped) != 2 {
		log(t, containerStopped, []string{"a", "c"})
		t.Fail()
	}
}

// Test when latest digest is pulled
func TestPullLatestDigest_Pulled(t *testing.T) {
	var pulledImage string

	listImages = func(_ *dockerClient.Client) []types.ImageSummary {
		return []types.ImageSummary{
			{
				RepoDigests: []string{"TestRepo@TestLatestDigest"},
			},
		}
	}

	pullImage = func(_ *dockerClient.Client, imageRef string) {
		pulledImage = imageRef
	}

	cake := Cake{
		LatestDigest: "TestLatestDigest",
		Repo:         "TestRepo",
	}

	cake.PullLatestDigest()

	if pulledImage != "" {
		log(t, pulledImage, "")
		t.Fail()
	}
}

// Test when latest digest is not pulled
func TestPullLatestDigest_NotPulled(t *testing.T) {
	var pulledImage string

	listImages = func(_ *dockerClient.Client) []types.ImageSummary {
		return []types.ImageSummary{}
	}

	pullImage = func(_ *dockerClient.Client, imageRef string) {
		pulledImage = imageRef
	}

	cake := Cake{
		LatestDigest: "TestLatestDigest",
		Repo:         "TestRepo",
	}

	cake.PullLatestDigest()

	if pulledImage != "TestRepo@TestLatestDigest" {
		log(t, pulledImage, "")
		t.Fail()
	}
}

// Test when latest digest is running but not under Cake's management
func TestRunLatestDigest_RunningNotControlled(t *testing.T) {
	expectedContainers := []string{"a", "c"}

	listContainers = func(client *dockerClient.Client) []types.Container {
		return []types.Container{
			{
				Image: "TestRepo@TestLatestDigest",
			},
		}
	}

	getRunningContainerIds = func(_ *dockerClient.Client, _ string, _ string) []string {
		return expectedContainers
	}

	cake := Cake{
		LatestDigest:      "TestLatestDigest",
		Repo:              "TestRepo",
		ContainersRunning: map[string]int{},
	}

	cake.RunLatestDigest()

	for _, id := range expectedContainers {
		if _, ok := cake.ContainersRunning[id]; !(ok) {
			t.Logf("Expected containers running %v, but Cake had %v", expectedContainers, cake.ContainersRunning)
			t.Fail()
		}
	}
}

// Test when latest container is running and is under cake's management
func TestRunLatestDigest_RunningAndControlled(t *testing.T) {
	expectedContainers := []string{"a", "c"}

	listContainers = func(client *dockerClient.Client) []types.Container {
		return []types.Container{
			{
				Image: "TestRepo@TestLatestDigest",
			},
		}
	}

	getRunningContainerIds = func(_ *dockerClient.Client, _ string, _ string) []string {
		return expectedContainers
	}

	cake := Cake{
		LatestDigest:      "TestLatestDigest",
		Repo:              "TestRepo",
		ContainersRunning: map[string]int{"a": 0, "c": 0},
	}

	cake.RunLatestDigest()

	for _, id := range expectedContainers {
		if _, ok := cake.ContainersRunning[id]; !(ok) {
			t.Logf("Expected containers running %v, but Cake had %v", expectedContainers, cake.ContainersRunning)
			t.Fail()
		}
	}
}

// Test when latest container is not running
func TestRunLatestDigest_NotRunningNotControlled(t *testing.T) {
	expectedContainer := "TestContainerId"

	listContainers = func(client *dockerClient.Client) []types.Container {
		return []types.Container{}
	}

	createContainer = func(_ *dockerClient.Client, _ container.Config, _ container.HostConfig, _ network.NetworkingConfig) string {
		return expectedContainer
	}

	cake := Cake{
		LatestDigest:      "TestLatestDigest",
		Repo:              "TestRepo",
		ContainersRunning: map[string]int{},
	}

	cake.RunLatestDigest()

	if _, ok := cake.ContainersRunning[expectedContainer]; !(ok) {
		log(t, cake.ContainersRunning, map[string]int{expectedContainer: 0})
		t.Fail()
	}
}
