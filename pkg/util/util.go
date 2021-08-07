// Package util contains utilities
package util

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

const (
	AppNameKey  = "appname"
	LogLevelKey = "logging.level"
	MusicKey    = "music"    // array of songs
	StationKey  = "stations" // stations to go with music
)

type Util struct {
	appName string
}

// SetupConfiguration is used to load/configure our configuration
func (util *Util) SetupConfiguration(appname, filename string) {
	viper.SetConfigName(filename)
	viper.SetConfigType("yaml")
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
	util.appName = appname
	viper.Set(AppNameKey, appname)
}

// SetupLogging is used to configure our logging once config is done
func (util *Util) SetupLogging() {
	desiredLogLevel := util.appName + "." + LogLevelKey
	logLevel, err := logrus.ParseLevel(viper.GetString(desiredLogLevel))
	if err != nil {
		logrus.Errorf("failed to parse log level: %s, and key: %s", viper.GetString(desiredLogLevel), desiredLogLevel)
		logLevel = logrus.ErrorLevel
	}
	logrus.Info("using a log level of: " + logLevel.String())
	logrus.SetLevel(logLevel)
}
