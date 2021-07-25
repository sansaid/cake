package main

import (
	"time"

	pb "github.com/sansaid/cake/pb"
)

// Run - run cake
func (c *Cake) Run(container *pb.Container) {
	c.GetLatestDigest(container)

	lastUpdateTime := time.Unix(container.LastUpdated, 0)
	lastCheckedTime := time.Unix(container.LastChecked, 0)

	if lastUpdateTime.After(lastCheckedTime) {
		c.StopPreviousDigest(container)
		c.PullLatestDigest(container)
		c.RunLatestDigest(container)
	}
}
