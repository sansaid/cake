package cmds

import (
	"sync"
)

type Config struct {
	mu    sync.Mutex
	image *string
}

// NewConfig - creates a new Config struct
func NewConfig(image *string) *Config {
	return &Config{
		image: image,
	}
}

func Run() {
	for {
		pollLatest()
	}
}
