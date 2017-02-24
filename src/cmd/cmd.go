// cmd/cmd.go
// Main entry point for gits

package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var configPath string

// MainCmd is the entry point for gits CLI
var MainCmd = &cobra.Command{
	Use:   "gits",
	Short: "a simple GitHub hook deployment application",
}

func init() {
	cobra.OnInitialize(initConfig)

	MainCmd.PersistentFlags().StringVarP(&configPath, "config", "c", "", "Configuration file to use")
}

func initConfig() {
	if configPath != "" {
		viper.SetConfigFile(configPath)
	}
}
