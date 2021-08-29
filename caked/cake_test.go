package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	pb "github.com/sansaid/cake/pb"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// func init(
// 	// disabling logging during tests
// 	log.SetOutput(ioutil.Discard)
// )

func TestListRunningContainerIds_Matches(t *testing.T) {
	mockDockerClient := new(MockDockerClient)
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

	mockHttpClient := new(MockHttpClient)
	cake := NewCakeWithMocks(mockDockerClient, mockHttpClient)

	expRes := []string{"Match1", "Match2"}
	res, err := cake.ListRunningContainerIds("testImage", "testDigest")

	assert.Equal(t, expRes, res)
	assert.Equal(t, nil, err)
}

func TestListRunningContainerIds_NoMatches(t *testing.T) {
	mockDockerClient := new(MockDockerClient)
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

	mockHttpClient := new(MockHttpClient)
	cake := NewCakeWithMocks(mockDockerClient, mockHttpClient)

	expRes := []string{}
	res, err := cake.ListRunningContainerIds("testImage", "testDigest")

	assert.Equal(t, expRes, res)
	assert.Equal(t, nil, err)
}

func TestListRunningContainerIds_Errors(t *testing.T) {
	mockDockerClient := new(MockDockerClient)
	mockDockerClient.On("ContainerList", mock.Anything, mock.Anything).Return(
		[]types.Container{},
		errors.New("test error"),
	)

	mockHttpClient := new(MockHttpClient)
	cake := NewCakeWithMocks(mockDockerClient, mockHttpClient)

	expRes := []string{}
	res, err := cake.ListRunningContainerIds("testImage", "testDigest")

	assert.Equal(t, expRes, res)
	assert.Error(t, err)
}

func TestCreateContainer_OK(t *testing.T) {
	pbContainer := &pb.Container{
		ImageName: "testImage@testDigest",
	}

	createdContainerId := "CreateMe"

	mockDockerClient := new(MockDockerClient)

	mockDockerClient.On("ContainerCreate", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, "").Return(
		container.ContainerCreateCreatedBody{
			ID: createdContainerId,
		},
		nil,
	)

	mockDockerClient.On("ContainerStart", mock.Anything, createdContainerId, types.ContainerStartOptions{}).Return(
		nil,
	)

	mockDockerClient.On("ContainerWait", mock.Anything, createdContainerId, container.WaitConditionNotRunning).Return(
		make(chan container.ContainerWaitOKBody),
		make(chan error),
	)

	mockHttpClient := new(MockHttpClient)
	cake := NewCakeWithMocks(mockDockerClient, mockHttpClient)

	expRes := createdContainerId
	res := cake.CreateContainer(pbContainer)

	assert.Equal(t, expRes, res)
}

func TestCreateContainer_ErrorOnCreate(t *testing.T) {
	pbContainer := &pb.Container{
		ImageName: "testImage@testDigest",
	}

	createdContainerId := "CreateMe"

	mockDockerClient := new(MockDockerClient)

	mockDockerClient.On("ContainerCreate", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, "").Return(
		container.ContainerCreateCreatedBody{},
		fmt.Errorf("test - error on create"),
	)

	mockDockerClient.On("ContainerStart", mock.Anything, createdContainerId, types.ContainerStartOptions{}).Return(
		nil,
	)

	mockDockerClient.On("ContainerWait", mock.Anything, createdContainerId, container.WaitConditionNotRunning).Return(
		make(chan container.ContainerWaitOKBody),
		make(chan error),
	)

	mockHttpClient := new(MockHttpClient)
	cake := NewCakeWithMocks(mockDockerClient, mockHttpClient)

	assert.PanicsWithValue(t, "Could not create container: test - error on create", func() { cake.CreateContainer(pbContainer) })
}

func TestCreateContainer_ErrorOnStart(t *testing.T) {
	pbContainer := &pb.Container{
		ImageName: "testImage@testDigest",
	}

	createdContainerId := "CreateMe"

	mockDockerClient := new(MockDockerClient)

	mockDockerClient.On("ContainerCreate", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, "").Return(
		container.ContainerCreateCreatedBody{
			ID: createdContainerId,
		},
		nil,
	)

	mockDockerClient.On("ContainerStart", mock.Anything, createdContainerId, types.ContainerStartOptions{}).Return(
		fmt.Errorf("test - error on start"),
	)

	mockDockerClient.On("ContainerWait", mock.Anything, createdContainerId, container.WaitConditionNotRunning).Return(
		make(chan container.ContainerWaitOKBody),
		make(chan error),
	)

	mockHttpClient := new(MockHttpClient)
	cake := NewCakeWithMocks(mockDockerClient, mockHttpClient)

	assert.PanicsWithValue(t, "Could not start container: test - error on start", func() { cake.CreateContainer(pbContainer) })
}

func TestGet_OK(t *testing.T) {
	mockDockerClient := new(MockDockerClient)
	mockHttpClient := new(MockHttpClient)

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

	cake := NewCakeWithMocks(mockDockerClient, mockHttpClient)

	cake.Get("test.url", jsonStruct)

	assert.Equal(t, "valA", jsonStruct.FieldA)
	assert.Equal(t, "valB", jsonStruct.FieldB)
}

func TestGet_ErrorOnResponse(t *testing.T) {
	mockDockerClient := new(MockDockerClient)
	mockHttpClient := new(MockHttpClient)

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

	cake := NewCakeWithMocks(mockDockerClient, mockHttpClient)

	assert.PanicsWithValue(t, "Could not read request response from test.url: test - error on response", func() { cake.Get("test.url", jsonStruct) })
}

func TestGet_MissingJsonFields(t *testing.T) {
	mockDockerClient := new(MockDockerClient)
	mockHttpClient := new(MockHttpClient)

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

	cake := NewCakeWithMocks(mockDockerClient, mockHttpClient)

	cake.Get("test.url", jsonStruct)

	assert.Equal(t, "", jsonStruct.FieldA)
	assert.Equal(t, "", jsonStruct.FieldB)
}

func TestGet_InvalidJson(t *testing.T) {
	mockDockerClient := new(MockDockerClient)
	mockHttpClient := new(MockHttpClient)

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

	cake := NewCakeWithMocks(mockDockerClient, mockHttpClient)

	assert.Panics(t, func() { cake.Get("test.url", jsonStruct) })
}

func TestGetLatestDigest_OK(t *testing.T) {
	mockDockerClient := new(MockDockerClient)
	mockHttpClient := new(MockHttpClient)
}
