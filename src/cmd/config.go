// cmd/config.go
// Configuration manipulation

package cmd

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// ConfigDefault holds a raw yaml format
var ConfigDefault = `---

log:
  level: info

server:
	port: 4200
`

// ConfigCmd holds all configuration manipulation commands
var ConfigCmd = &cobra.Command{
	Use: "config",
}

// ConfigViewCmd prints out the current configuration file
var ConfigViewCmd = &cobra.Command{
	Use:   "view",
	Short: "prints out configuration file",
	Run: func(cmd *cobra.Command, args []string) {
		viper.ReadInConfig()
		config := viper.ConfigFileUsed()
		if config == "" {
			return
		}

		b, err := ioutil.ReadFile(config)
		if err != nil {
			fmt.Fprint(os.Stderr, "Unable to read configuration file\n")
			fmt.Fprintf(os.Stderr, "%s\n", err.Error())
			os.Exit(1)
		}

		fmt.Print(string(b))
		os.Exit(0)
	},
}

// ConfigGenerateCmd generates a default configuration file
var ConfigGenerateCmd = &cobra.Command{
	Use:   "generate [path]",
	Short: "creates a default configuration file",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			fmt.Fprint(os.Stderr, "we need a path to generate configuration at\n")
			os.Exit(1)
		}

		b := []byte(ConfigDefault)

		if err := ioutil.WriteFile(args[0], b, 0644); err != nil {
			fmt.Fprint(os.Stderr, "Unable to write configuration file\n")
			fmt.Fprintf(os.Stderr, "%s\n", err.Error())
			os.Exit(1)
		}

		os.Exit(0)
	},
}

func init() {
	ConfigCmd.AddCommand(ConfigViewCmd)
	ConfigCmd.AddCommand(ConfigGenerateCmd)

	MainCmd.AddCommand(ConfigCmd)
}
