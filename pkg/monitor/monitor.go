// Package monitor is the main monitoring package
package monitor

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	"github.com/milligan22963/radio/pkg/util"
	"github.com/sirupsen/logrus"
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
		fmt.Println("signal detected, ending operation")
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

func (monitor *MonitorCmd) loadStatic(utilities util.Util) {
	monitor.mixer.Initialize()

	for _, v := range utilities.RadioInformation.Static {
		logrus.Infof("loading static file: %s", v)
		err := monitor.mixer.AddStatic(v)
		if err != nil {
			logrus.Error(err)
		}
	}
}

func (monitor *MonitorCmd) loadStations(utilities util.Util) {
	for _, v := range utilities.RadioInformation.Stations {
		logrus.Infof("adding station: %f - file: %s", v.Station, v.Music)
		err := monitor.mixer.AddStation(float64(v.Station), v.Music)
		if err != nil {
			logrus.Errorf("unable to process station: %v", v)
			continue
		}
	}
}

func (monitor *MonitorCmd) setup() {
	// configure gpios

}

func (monitor *MonitorCmd) Run(utilities util.Util) {

	monitor.loadStatic(utilities)
	monitor.loadStations(utilities)

	monitor.setup()

	go monitor.playRadio()

	monitor.waitForExit()

	monitor.mixer.Close()
}
