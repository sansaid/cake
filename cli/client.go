package main

import (
	"fmt"
	"os"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var CfgFile string
var Image string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "cake",
	Short: "Container manager for tiny servers",
	Long: `Cake manages your running containers on your tiny home server (e.g. Raspberry Pi).
It polls a container registry of your choice and automatically updates the running container on your tiny
server. You'll never have to login to your tiny server to update your running code.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { fmt.Println("Hello this is cake") },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVarP(&CfgFile, "config", "c", "", "config file (default is $HOME/.cake.yaml)")
	rootCmd.PersistentFlags().StringVarP(&Image, "image", "i", "", "[required] name of the image to manage (e.g. repo/image:tag) - cake currently only supports images in Docker Hub")

	rootCmd.MarkPersistentFlagRequired("image")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if CfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(CfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".cake" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".cake")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}

func main() {
	Execute()
}
