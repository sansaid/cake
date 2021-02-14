package cmds

import (
	"sync"
)

type Config struct {
	mu          sync.Mutex
	registry    string            // The name of the registry to poll
	imageTagMap map[string]string // The map of images to their tags that will be polled for changes
}

// NewConfig - creates a new Config struct
func NewConfig(registry string) *Config {
	return &Config{
		registry:    registry,
		imageTagMap: make(map[string]string),
	}
}

// AddImageTagRef - adds an image and its associated tag to poll
func (c *Config) AddImageTagRef(image string, tag string) *Config {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.imageTagMap[image] = tag

	return c
}
