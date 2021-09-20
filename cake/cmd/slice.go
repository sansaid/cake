package cmd

import (
	"github.com/sansaid/cake/cake/pb"
)

func NewSlice(imageName string) *pb.Slice {
	os := "linux"
	arch := "amd64"

	return &pb.Slice{
		ImageName:    imageName,
		Os:           os,
		Architecture: arch,
	}
}
