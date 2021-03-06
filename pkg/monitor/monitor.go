// Package monitor is the main monitoring package
package monitor

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
	"github.com/milligan22963/radio/pkg/util"
	"github.com/sirupsen/logrus"
	"github.com/stianeikeland/go-rpio/v4"
)

var gpioSleepPeriod = time.Second / 2

type MonitorCmd struct {
	mixer           util.RadioMixer
	applicationDone chan bool
}

func (monitor *MonitorCmd) waitForExit() {
	signals := make(chan os.Signal, 1)
	monitor.applicationDone = make(chan bool, 1)

	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-signals
		fmt.Println("signal detected, ending operation")
		fmt.Println(sig)
		monitor.applicationDone <- true
	}()

	<-monitor.applicationDone
}

func (monitor *MonitorCmd) playRadio() {
	logrus.Info("starting to play the radio")
	desiredFormat := beep.SampleRate(44100)
	err := speaker.Init(desiredFormat, desiredFormat.N(time.Second/10))
	if err != nil {
		logrus.Errorf("speaker init fail: %v\n", err)
	}
	speaker.Play(beep.Seq(&monitor.mixer, beep.Callback(func() {
		logrus.Info("song has ended.")
		monitor.applicationDone <- true
	})))

	<-monitor.applicationDone
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

func (monitor *MonitorCmd) MonitorGPIO(station float64, pin rpio.Pin) {
	pin.Input()
	pin.PullUp()
	pin.Detect(rpio.FallEdge) // enable falling edge event detection

	defer pin.Detect(rpio.NoEdge) // disable edge event detection

	for {
		select {
		case <-monitor.applicationDone:
			return
		default:
			if pin.EdgeDetected() { // check if event occured
				monitor.mixer.PlayStation(station)
			}
			time.Sleep(gpioSleepPeriod)
		}
	}
}

func (monitor *MonitorCmd) setupGPIO(utilities util.Util) {
	// configure gpios
	for _, v := range utilities.RadioInformation.Stations {
		logrus.Infof("configuring gpio: %f - pin: %d", v.Station, v.GPIO)

		go monitor.MonitorGPIO(float64(v.Station), rpio.Pin(v.GPIO))
	}
}

func (monitor *MonitorCmd) Run(utilities util.Util) {

	err := rpio.Open()

	if err != nil {
		logrus.Errorf("unable to open gpios: %v", err)
		return
	}

	defer rpio.Close()

	monitor.loadStatic(utilities)
	monitor.loadStations(utilities)

	monitor.setupGPIO(utilities)

	go monitor.playRadio()

	monitor.waitForExit()

	monitor.mixer.Close()
}
