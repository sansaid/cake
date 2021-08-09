package main

import (
	"net/http"
	"sync"
	"time"

	dockerClient "github.com/docker/docker/client"
	pb "github.com/sansaid/cake/pb"
	"github.com/sansaid/cake/utils"
)

type ContainerList struct {
	sync.RWMutex
	containers map[string]struct{}
}

type Cake struct {
	pb.UnimplementedCakedServer
	DockerClient      ModContainerAPIClient
	HttpClient        ModHTTPClient
	ContainersRunning ContainerList
	StopTimeout       time.Duration
}

func NewCake() *Cake {
	client, err := dockerClient.NewEnvClient()

	utils.Check(err, "Cannot create cake client")

	return &Cake{
		DockerClient: client,
		HttpClient:   &http.Client{Timeout: time.Second * 3},
		StopTimeout:  30 * time.Second,
	}
}
