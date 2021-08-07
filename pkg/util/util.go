// Package util contains utilities
package util

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func setupConfiguration(appname, filename string) {
	viper.SetConfigName(filename)
	viper.AddConfigPath("/etc/" + appname + "/")
	viper.AddConfigPath("$HOME/." + appname)
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Failed to read in the config file
			logrus.Info("using configuration file: " + filename)
		} else {
			logrus.Error("unable to read in configuration file: " + filename)
		}
	}
}
