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
	"github.com/stretchr/testify/mock"
)

// ------- MockDockerClient -------
type MockHttpClient struct {
	mock.Mock
}

func (m *MockHttpClient) Do(req *http.Request) (*http.Response, error) {
	args := m.Called(req)

	return args.Get(0).(*http.Response), args.Error(1)
}

// ------- MockDockerClient -------
type MockDockerClient struct {
	mock.Mock
}

func (m *MockDockerClient) ContainerList(ctx context.Context, options types.ContainerListOptions) ([]types.Container, error) {
	args := m.Called(ctx, options)

	return args.Get(0).([]types.Container), args.Error(1)
}

func (m *MockDockerClient) ImagePull(ctx context.Context, refStr string, options types.ImagePullOptions) (io.ReadCloser, error) {
	args := m.Called(ctx, refStr, options)

	return args.Get(0).(io.ReadCloser), args.Error(1)
}

func (m *MockDockerClient) ImageList(ctx context.Context, options types.ImageListOptions) ([]types.ImageSummary, error) {
	args := m.Called(ctx, options)

	return args.Get(0).([]types.ImageSummary), args.Error(1)
}

func (m *MockDockerClient) ContainerCreate(ctx context.Context, config *container.Config, hostConfig *container.HostConfig, networkingConfig *network.NetworkingConfig, platform *v1.Platform, containerName string) (container.ContainerCreateCreatedBody, error) {
	args := m.Called(ctx, config, hostConfig, networkingConfig, platform, containerName)

	return args.Get(0).(container.ContainerCreateCreatedBody), args.Error(1)
}

func (m *MockDockerClient) ContainerStart(ctx context.Context, containerID string, options types.ContainerStartOptions) error {
	args := m.Called(ctx, containerID, options)

	return args.Error(0)
}

func (m *MockDockerClient) ContainerWait(ctx context.Context, containerID string, condition container.WaitCondition) (<-chan container.ContainerWaitOKBody, <-chan error) {
	args := m.Called(ctx, containerID, condition)

	return args.Get(0).(<-chan container.ContainerWaitOKBody), args.Get(1).(<-chan error)
}

func (m *MockDockerClient) ContainerStop(ctx context.Context, containerID string, timeout *time.Duration) error {
	args := m.Called(ctx, containerID, timeout)

	return args.Error(0)
}

func (m *MockDockerClient) Close() error {
	args := m.Called()

	return args.Error(0)
}

// ------- MockCake -------
func NewMockCake() *Cake {
	return &Cake{
		DockerClient: new(MockDockerClient),
		StopTimeout:  30 * time.Second,
		HttpClient:   new(MockHttpClient),
	}
}
