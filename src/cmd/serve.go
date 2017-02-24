// cmd/serve.go
// Starts gits web server

package cmd

import (
	"fmt"
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/btkostner/gits/src/config"
	"github.com/btkostner/gits/src/controller"
	"github.com/btkostner/gits/src/server"
)

// ServeCmd starts the web server for listening to GitHub hooks
var ServeCmd = &cobra.Command{
	Use:   "serve",
	Short: "starts the web server to listen for hooks",
	Run: func(cmd *cobra.Command, args []string) {
		if err := config.Setup(); err != nil {
			fmt.Fprint(os.Stderr, "Unable to read configuration\n")
			fmt.Fprintf(os.Stderr, "%s\n", err.Error())
			os.Exit(1)
		}

		p := viper.GetInt("server.port")
		if p == 0 {
			fmt.Fprint(os.Stderr, "Server port must be set\n")
			os.Exit(1)
		}

		s := server.NewServer(controller.New())

		logrus.Infof("Starting web server on port %v", p)
		if err := s.Serve(p); err != nil {
			logrus.Error(err)
			os.Exit(1)
		}

		os.Exit(0)
	},
}

func init() {
	var p int

	ServeCmd.PersistentFlags().IntVarP(&p, "port", "p", 4200, "Port to listen on")

	viper.BindPFlag("server.port", ServeCmd.PersistentFlags().Lookup("port"))

	MainCmd.AddCommand(ServeCmd)
}
