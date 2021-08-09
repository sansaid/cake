package main

import (
	"context"
	"io"
	"net/http"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	v1 "github.com/opencontainers/image-spec/specs-go/v1"
)

// ------- MockReadCloser -------
type MockReadCloser struct{}

func (m *MockReadCloser) Read(p []byte) (n int, err error) {
	return 0, nil
}

func (m *MockReadCloser) Close() error {
	return nil
}

// ------- MockDockerClient -------
type MockHttpClient struct{}

func (m *MockHttpClient) Do(req *http.Request) (*http.Response, error) {
	return http.Response{}, nil
}

// ------- MockDockerClient -------
type MockDockerClient struct{}

func (m *MockDockerClient) ContainerList(ctx context.Context, options types.ContainerListOptions) ([]types.Container, error) {
	return []types.Container{}, nil
}

func (m *MockDockerClient) ImagePull(ctx context.Context, refStr string, options types.ImagePullOptions) (io.ReadCloser, error) {
	return &MockReadCloser{}, nil
}

func (m *MockDockerClient) ImageList(ctx context.Context, options types.ImageListOptions) ([]types.ImageSummary, error) {
	return []types.ImageSummary{}, nil
}

func (m *MockDockerClient) ContainerCreate(ctx context.Context, config *container.Config, hostConfig *container.HostConfig, networkingConfig *network.NetworkingConfig, platform *v1.Platform, containerName string) (container.ContainerCreateCreatedBody, error) {
	return container.ContainerCreateCreatedBody{}, nil
}

func (m *MockDockerClient) ContainerStart(ctx context.Context, containerID string, options types.ContainerStartOptions) error {
	return nil
}

func (m *MockDockerClient) ContainerWait(ctx context.Context, containerID string, condition container.WaitCondition) (<-chan container.ContainerWaitOKBody, <-chan error) {
	return make(<-chan container.ContainerWaitOKBody), make(<-chan error)
}

// ------- MockCake -------
func NewMockCake() *Cake {
	return &Cake{
		DockerClient: &MockDockerClient{},
		StopTimeout:  30 * time.Second,
		HttpClient:   &MockHttpClient{},
	}
}
