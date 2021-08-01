// Package monitor is the main monitoring package
package monitor

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	"github.com/sirupsen/logrus"
)

type MonitorCmd struct {
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

func (monitor *MonitorCmd) playMusicFile(file string) {
	f, err := os.Open(file)
	if err != nil {
		logrus.Errorf("failed to open file: %v", err)
	}

	streamer, format, err := mp3.Decode(f)
	if err != nil {
		logrus.Errorf("mp3 decoding fail: %v", err)
	}
	defer streamer.Close()
	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	speaker.Play(streamer)
}

func (monitor *MonitorCmd) setup() {
	// configure gpios
}

func (monitor *MonitorCmd) WatchGPIOS() {
	// on gpio play sound otherwise static
}

func (monitor *MonitorCmd) Run() {

	go func() {
		monitor.WatchGPIOS()
	}()

	monitor.waitForExit()
}
