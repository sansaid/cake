package cmd

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/docker/docker/api/types"
	cakeMocks "github.com/sansaid/cake/caked/mocks"
	"github.com/sansaid/cake/caked/pb"
	"github.com/sansaid/cake/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAddNewSlice_OK(t *testing.T) {
	var wg *sync.WaitGroup
	_, cancel := context.WithCancel(context.Background())

	slices := []*pb.Slice{
		{ImageName: "sansaid/testimageuno:1234"},
		{ImageName: "sansaid/testimagedos:4567"},
	}

	cake := NewCake()

	for _, slice := range slices {
		cake.AddNewSlice(slice, wg, cancel)
	}

	assert.Len(t, cake.RunningSlices.slices, 2)
	assert.Contains(t, cake.RunningSlices.slices, "sansaid/testimageuno:1234", "sansaid/testimagedos:4567")
}

func TestMarshallHttp_OK(t *testing.T) {
	mockHttpClient := new(cakeMocks.CakeHTTPClient)

	req, _ := http.NewRequest(http.MethodGet, "test.url", nil)

	mockHttpClient.On("Do", req).Return(
		&http.Response{
			Body: io.NopCloser(bytes.NewBufferString("{ \"fieldA\": \"valA\", \"fieldB\": \"valB\" }")),
		},
		nil,
	)

	jsonStruct := &struct {
		FieldA string `json:"fieldA"`
		FieldB string `json:"fieldB"`
	}{}

	cake := NewCake(WithHttpClient(mockHttpClient))

	cake.MarshalHttp("test.url", jsonStruct)

	assert.Equal(t, "valA", jsonStruct.FieldA)
	assert.Equal(t, "valB", jsonStruct.FieldB)
}

func TestMarshallHttp_ErrorOnResponse(t *testing.T) {
	mockHttpClient := new(cakeMocks.CakeHTTPClient)

	req, _ := http.NewRequest(http.MethodGet, "test.url", nil)

	mockHttpClient.On("Do", req).Return(
		&http.Response{
			Body: io.NopCloser(bytes.NewBufferString("{ \"fieldA\": \"valA\", \"fieldB\": \"valB\" }")),
		},
		fmt.Errorf("test - error on response"),
	)

	jsonStruct := &struct {
		FieldA string `json:"fieldA"`
		FieldB string `json:"fieldB"`
	}{}

	cake := NewCake(WithHttpClient(mockHttpClient))

	assert.Error(t, cake.MarshalHttp("test.url", jsonStruct))
}

func TestMarshallHttp_MissingJsonFields(t *testing.T) {
	mockHttpClient := new(cakeMocks.CakeHTTPClient)

	req, _ := http.NewRequest(http.MethodGet, "test.url", nil)

	mockHttpClient.On("Do", req).Return(
		&http.Response{
			Body: io.NopCloser(bytes.NewBufferString("{ \"fieldC\": \"valA\", \"fieldD\": \"valB\" }")),
		},
		nil,
	)

	jsonStruct := &struct {
		FieldA string `json:"fieldA"`
		FieldB string `json:"fieldB"`
	}{}

	cake := NewCake(WithHttpClient(mockHttpClient))

	err := cake.MarshalHttp("test.url", jsonStruct)

	assert.Equal(t, "", jsonStruct.FieldA)
	assert.Equal(t, "", jsonStruct.FieldB)
	assert.NoError(t, err)
}

func TestMarshallHttp_InvalidJson(t *testing.T) {
	mockHttpClient := new(cakeMocks.CakeHTTPClient)

	req, _ := http.NewRequest(http.MethodGet, "test.url", nil)

	mockHttpClient.On("Do", req).Return(
		&http.Response{
			Body: io.NopCloser(bytes.NewBufferString("bad json")),
		},
		nil,
	)

	jsonStruct := &struct {
		FieldA string `json:"fieldA"`
		FieldB string `json:"fieldB"`
	}{}

	cake := NewCake(WithHttpClient(mockHttpClient))

	assert.Error(t, cake.MarshalHttp("test.url", jsonStruct))
}

func TestListRunningContainerIds_Matches(t *testing.T) {
	mockDockerClient := new(cakeMocks.CakeContainerAPIClient)
	mockDockerClient.On("ContainerList", mock.Anything, mock.Anything).Return(
		[]types.Container{
			{
				ID:    "Match1",
				Image: "testImage@testDigest",
			},
			{
				ID:    "Match2",
				Image: "testImage@testDigest",
			},
			{
				ID:    "NotMatch",
				Image: "testImage@notTestDigest",
			},
		},
		nil,
	)

	cake := NewCake(WithDockerClient(mockDockerClient))

	expRes := []string{"Match1", "Match2"}
	res, err := cake.ListRunningContainerIds("testImage", "testDigest")

	assert.Equal(t, expRes, res)
	assert.Equal(t, nil, err)
}

func TestListRunningContainerIds_NoMatches(t *testing.T) {
	mockDockerClient := new(cakeMocks.CakeContainerAPIClient)
	mockDockerClient.On("ContainerList", mock.Anything, mock.Anything).Return(
		[]types.Container{
			{
				ID:    "NotMatch1",
				Image: "testImage@notTestDigest",
			},
			{
				ID:    "NotMatch2",
				Image: "testImage@notTestDigest",
			},
		},
		nil,
	)

	cake := NewCake(WithDockerClient(mockDockerClient))

	expRes := []string{}
	res, err := cake.ListRunningContainerIds("testImage", "testDigest")

	assert.Equal(t, expRes, res)
	assert.Equal(t, nil, err)
}

func TestListRunningContainerIds_Errors(t *testing.T) {
	mockDockerClient := new(cakeMocks.CakeContainerAPIClient)
	mockDockerClient.On("ContainerList", mock.Anything, mock.Anything).Return(
		[]types.Container{},
		errors.New("test error"),
	)

	cake := NewCake(WithDockerClient(mockDockerClient))

	expRes := []string{}
	res, err := cake.ListRunningContainerIds("testImage", "testDigest")

	assert.Equal(t, expRes, res)
	assert.Error(t, err)
}

// TODO - NEXT: finish writing tests for GetLatestDigest ->
// Need to read from fixtures for OK tests
// Need to write test for when MarshallHttp returns error
func TestGetLatestDigest_OK(t *testing.T) {
	mockHttpClient := new(cakeMocks.CakeHTTPClient)
	imageName := "sansaid/dummyimage:dummytag"
	repoUrl := fmt.Sprintf("https://registry.hub.docker.com/v2/repositories/%s/tags?ordering=last_updated", imageName)

	req, _ := http.NewRequest(http.MethodGet, repoUrl, nil)

	resp := RepoList{
		Count: 2,
		Results: []ImageDetails{
			{
				ID:      1,
				ImageID: "dummyImageIdUno",
				Images: []Image{
					{
						Architecture: "amd64",
						OS:           "linux",
						Digest:       "dummyDigestUno",
						LastPushed:   time.Date(2020, time.January, 4, 1, 0, 0, 0, time.UTC),
					},
				},
				TagLastPushed: time.Date(2020, time.January, 4, 0, 0, 0, 0, time.UTC),
			},
			{
				ID:      2,
				ImageID: "dummyImageIdDos",
				Images: []Image{
					{
						Architecture: "amd64",
						OS:           "linux",
						Digest:       "dummyDigestDos",
						LastPushed:   time.Date(2020, time.January, 3, 1, 0, 0, 0, time.UTC),
					},
				},
				TagLastPushed: time.Date(2020, time.January, 3, 0, 0, 0, 0, time.UTC),
			},
		},
	}

	mockHttpClient.On("Do", req).Return(
		&http.Response{
			Body: utils.JsonNopCloser(resp),
		},
		nil,
	)

	cake := NewCake(WithHttpClient(mockHttpClient))

	expDigest := "dummyDigestUno"
	expTime := time.Date(2020, time.January, 4, 1, 0, 0, 0, time.UTC).Unix()

	digest, time, err := cake.GetLatestDigest(&pb.Slice{
		ImageName:    imageName,
		Os:           "linux",
		Architecture: "amd64",
	})

	assert.Equal(t, expDigest, digest)
	assert.Equal(t, expTime, time)
	assert.NoError(t, err)
}

// func TestGetLatestDigest_OK(t *testing.T) {
// 	mockHttpClient := new(cakeMocks.CakeHTTPClient)

// 	req, _ := http.NewRequest(http.MethodGet, "test.url", nil)

// 	mockHttpClient.On("Do", req).Return(
// 		&http.Response{
// 			Body: io.NopCloser(bytes.NewBufferString("{ \"fieldA\": \"valA\", \"fieldB\": \"valB\" }")),
// 		},
// 		nil,
// 	)

// 	cake := NewCake(WithHttpClient(mockHttpClient))

// 	expRes := []string{}
// 	res, err := cake.ListRunningContainerIds("testImage", "testDigest")

// 	assert.Equal(t, expRes, res)
// 	assert.Error(t, err)
// }
