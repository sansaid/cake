/*
Copyright © 2021 SANYIA SAIDOVA

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
	"context"

	pb "github.com/sansaid/cake/pb"
	"github.com/sansaid/cake/utils"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

// stopCmd represents the stop command
var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop a caked image",
	Long:  `Stop containers associated with this image. The container will no longer run and you will no longer receive updates from Docker Hub.`,
	Run: func(cmd *cobra.Command, args []string) {
		var opts []grpc.DialOption

		conn, err := grpc.Dial("localhost:6010", opts...)
		defer conn.Close()

		utils.Check(err, "Cannot initialise gRPC dial")

		container := NewContainer(Image)
		client := pb.NewCakedClient(conn)

		client.StopContainer(context.Background(), container)
	},
}

func init() {
	rootCmd.AddCommand(stopCmd)
}
