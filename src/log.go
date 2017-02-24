// log.go
// Sets up logrus log level

package main

import (
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

func updateLogLevel() {
	switch strings.ToUpper(viper.GetString("log.level")) {
	case "DEBUG":
		logrus.SetLevel(logrus.DebugLevel)
	case "INFO":
		logrus.SetLevel(logrus.InfoLevel)
	case "WARN":
		logrus.SetLevel(logrus.WarnLevel)
	case "ERROR":
		logrus.SetLevel(logrus.ErrorLevel)
	default:
		logrus.SetLevel(logrus.InfoLevel)
	}

}

func init() {
	viper.OnConfigChange(func(e fsnotify.Event) {
		updateLogLevel()
	})

	updateLogLevel()
}
