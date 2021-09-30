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
	"context"

	"github.com/sansaid/cake/cakectl/pb"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

// stopCmd represents the stop command
var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop running your image as a Cake slice",
	Long:  `Stop running your image as a Cake slice. This will kill any containers running with this image.`,
	Run: func(cmd *cobra.Command, args []string) {
		var opts []grpc.DialOption
		conn, err := grpc.Dial("localhost:6010", opts...) // TODO: these should be CLI arguments or config
		defer conn.Close()

		image := &pb.Image{
			Name: sliceImage,
		}

		if err != nil {
			log.Errorf("Could not stop slice for image %s: %s", image.Name, err) // TODO: decide if we need to also include stack in some of the error messages
			return
		}

		client := pb.NewCakedClient(conn)

		status, err := client.StopSlice(context.Background(), image)

		if err != nil || status.Status != 0 { // TODO: Remove reliance on SliceStatus - rely only on error message
			log.Errorf("Failed to run slice for image %s: %s - %s", image.Name, status.Message, err) // TODO: decide if we need to also include stack in some of the error messages
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(stopCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// stopCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// stopCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
