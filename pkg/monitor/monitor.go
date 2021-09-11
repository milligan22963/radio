// Package monitor is the main monitoring package
package monitor

import (
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	"github.com/milligan22963/radio/pkg/util"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type MonitorCmd struct {
	mixer util.RadioMixer
}

func (monitor *MonitorCmd) waitForExit() {
	signals := make(chan os.Signal, 1)
	doneFlag := make(chan bool, 1)

	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-signals
		fmt.Println("ending operation")
		fmt.Println(sig)
		doneFlag <- true
	}()

	<-doneFlag
}

func (monitor *MonitorCmd) playRadio() {
	logrus.Info("starting to play the radio")
	desiredFormat := beep.SampleRate(44100)
	err := speaker.Init(desiredFormat, desiredFormat.N(time.Second/10))
	if err != nil {
		logrus.Errorf("speaker init fail: %v\n", err)
	}
	done := make(chan bool)
	speaker.Play(beep.Seq(&monitor.mixer, beep.Callback(func() {
		logrus.Info("song has ended.")
		done <- true
	})))

	<-done
}

func (monitor *MonitorCmd) playMusicFile(file string) {
	f, err := os.Open(file)
	if err != nil {
		logrus.Errorf("failed to open file: %v\n", err)
	}

	streamer, format, err := mp3.Decode(f)
	if err != nil {
		logrus.Errorf("mp3 decoding fail: %v\n", err)
	}
	defer streamer.Close()
	err = speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	if err != nil {
		logrus.Errorf("speaker init fail: %v\n", err)
	}
	done := make(chan bool)
	speaker.Play(beep.Seq(streamer, beep.Callback(func() {
		logrus.Info("song has ended.")
		done <- true
	})))

	<-done
}

func (monitor *MonitorCmd) loadStatic() {
	appKey := viper.GetString(util.AppNameKey)

	staticKey := appKey + "." + util.StaticKey
	staticFiles := viper.GetStringSlice(staticKey)

	monitor.mixer.Initialize()

	for _, v := range staticFiles {
		err := monitor.mixer.AddStatic(v)
		if err != nil {
			logrus.Error(err)
		}
	}
}

func (monitor *MonitorCmd) loadStations() {
	appKey := viper.GetString(util.AppNameKey)

	songKey := appKey + "." + util.MusicKey
	songs := viper.GetStringSlice(songKey)

	stationKey := appKey + "." + util.StationKey
	stations := viper.GetStringSlice(stationKey)

	if len(songs) != len(stations) {
		logrus.Warnf("songs:\n+%v", songs)
		logrus.Warnf("stations:\n+%v", stations)
		panic("mismatch between number of songs and stations")
	}
	for k, v := range songs {
		station, err := strconv.ParseFloat(stations[k], 10)
		if err != nil {
			logrus.Errorf("unable to parse station: %v, skipping", stations[k])
			continue
		}
		err = monitor.mixer.AddStation(station, v)
		if err != nil {
			logrus.Errorf("unable to process station: %v", v)
			continue
		}
	}
}

func (monitor *MonitorCmd) setup() {
	// configure gpios
}

func (monitor *MonitorCmd) Run() {

	monitor.loadStatic()
	monitor.loadStations()

	monitor.setup()

	monitor.playRadio()

	monitor.waitForExit()

	monitor.mixer.Close()
}
