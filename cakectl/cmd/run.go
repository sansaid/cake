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

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run your image as a Cake slice",
	Long: `Run your image as a Cake slice. A Cake slice is simply a Docker container that is always kept up to \
	date with an image in your Docker Hub registry. Only compatible with public Docker Hub registries at the moment.`,
	Run: func(cmd *cobra.Command, args []string) {
		var opts []grpc.DialOption
		conn, err := grpc.Dial("localhost:6010", opts...) // TODO: these should be CLI arguments or config
		defer conn.Close()

		if err != nil {
			log.Errorf("Could not create slice for image %s: %s", sliceImage, err) // TODO: decide if we need to also include stack in some of the error messages
			return
		}

		client := pb.NewCakedClient(conn)
		slice := NewSlice(sliceImage)

		status, err := client.RunSlice(context.Background(), slice)

		if err != nil || status.Status != 0 {
			log.Errorf("Failed to run slice for image %s: %s - %s", sliceImage, status.Message, err) // TODO: decide if we need to also include stack in some of the error messages
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(runCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// runCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// runCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
