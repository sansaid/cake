package main

import (
	"fmt"
	"os"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// TODO: confirm the right ports to use - https://www.iana.org/assignments/service-names-port-numbers/service-names-port-numbers.xhtml?search=6010
const port = 6010

var Registry string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "caked",
	Short: "Cake daemon manager",
	Long:  `The Cake daemon`,
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
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	// Find home directory.
	home, err := homedir.Dir()
	cobra.CheckErr(err)

	// Search config in home directory with name ".caked" (without extension).
	viper.AddConfigPath(home)
	viper.SetConfigName(".caked")

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}

func main() {
	Execute()
}
