/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"strings"

	pb "github.com/sansaid/cake/pb"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

const registry string = "https://hub.docker.com"

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start a caked image",
	Long:  `Start running an image managed by cake. Any changes to the image in Docker Hub will automatically get updated where the corresponding image is running.`,
	Run: func(cmd *cobra.Command, args []string) {
		splitImage := strings.Split(imageName, ":")
		repo, tag := splitImage[0], splitImage[1]

		CreateContainer()

		var opts []grpc.DialOption

		conn, err := grpc.Dial("localhost:6010", opts...)

		if err != nil {
			panic(err)
		}

		defer conn.Close()

		client := pb.NewCakedClient(conn)

		client.StartContainer()
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
}
