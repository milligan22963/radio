// Package util contains utilities
package util

import (
	"io/ioutil"
	"path/filepath"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

const (
	AppNameKey  = "appname"
	LogLevelKey = "logging.level"
	MusicKey    = "music"    // song
	StationsKey = "stations" // stations to go with music
	StationKey  = "station"  // single station
	StaticKey   = "static"   // static options
	GPIOKey     = "gpio"     // the gpio to use
)

type MusicStation struct {
	Music   string  `yaml:"music"`
	Station float32 `yaml:"station"`
	GPIO    int     `yaml:"gpio"`
}

type RadioInformation struct {
	Static   []string       `yaml:"static"`
	Stations []MusicStation `yaml:"stations"`
}

type Logging struct {
	LogLevel string `yaml:"level"`
}

// Util is a struct tracking the app nam
type Util struct {
	RadioInformation RadioInformation `yaml:"radio"`
	Logging          Logging          `yaml:"logging"`
}

// SetupConfiguration is used to load/configure our configuration
func (util *Util) SetupConfiguration(fileName string) {
	fileContents, err := ioutil.ReadFile(filepath.Clean(fileName))

	if err != nil {
		logrus.Errorf("unable to load radio station information: %v", err)
		return
	}

	err = yaml.Unmarshal(fileContents, util)
	if err != nil {
		logrus.Errorf("unable to load radio station information: %v", err)
		return
	}

	logrus.Infof("configuration: %v", util)
}

// SetupLogging is used to configure our logging once config is done
func (util *Util) SetupLogging() {
	logLevel, err := logrus.ParseLevel(util.Logging.LogLevel)
	if err != nil {
		logrus.Errorf("failed to parse log level: %s", util.Logging.LogLevel)
		logLevel = logrus.ErrorLevel
	}
	logrus.Info("using a log level of: " + logLevel.String())
	logrus.SetLevel(logLevel)
}
