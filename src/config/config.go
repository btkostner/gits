// config/config.go
// Sets up viper configuration loading

package config

import (
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

func init() {
	home := os.Getenv("HOME")
	xdg := os.Getenv("XDG_CONFIG_HOME")

	viper.SetConfigName("config")
	viper.SetConfigType("yml")

	viper.AddConfigPath("/usr/share/gits/")
	viper.AddConfigPath("/etc/gits/")
	if xdg != "" {
		viper.AddConfigPath(xdg + "/gits/")
	} else {
		viper.AddConfigPath(home + "/.local/share/gits")
	}
	viper.AddConfigPath(".")

	viper.SetEnvPrefix("gits")
	viper.AutomaticEnv()

	viper.SetDefault("log.level", "info")
	viper.SetDefault("server.port", 4200)

	viper.WatchConfig()
	Setup()
}

// Setup reads and sets up Config settings
func Setup() error {
	viper.OnConfigChange(func(e fsnotify.Event) {
		logrus.Debug("Configuration updated")
		ReadProjects()
	})

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	ReadProjects()
	return nil
}
