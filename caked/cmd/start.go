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
	"fmt"
	"net"

	"github.com/sansaid/cake/caked/pb"
	"github.com/sansaid/cake/utils"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

var cake *Cake

func startCake(port int) {
	lis := utils.Must(net.Listen("tcp", fmt.Sprintf("localhost:%d", port)))
	defer lis.(net.Listener).Close()

	var opts []grpc.ServerOption

	grpcServer := grpc.NewServer(opts...)

	cake = NewCake()

	pb.RegisterCakedServer(grpcServer, cake)

	grpcServer.Serve(lis.(net.Listener))
}

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the Cake daemon",
	Long:  `Start running the Cake daemon`,
	Run: func(cmd *cobra.Command, args []string) {
		startCake(6010) // TODO: turn into server config: --port
	},
}

func init() {
	rootCmd.AddCommand(startCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// startCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// startCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
