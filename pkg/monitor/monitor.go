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
	"github.com/spf13/viper"
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/platforms/keyboard"
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

func (monitor *MonitorCmd) setup() {

	// configure gpios
	keys := keyboard.NewDriver()

	appKey := viper.GetString(util.AppNameKey)

	songKey := appKey + "." + util.MusicKey
	songs := viper.GetStringSlice(songKey)

	work := func() {
		keys.On(keyboard.Key, func(data interface{}) {
			key := data.(keyboard.KeyEvent)

			if key.Key == keyboard.A {
				logrus.Info("start to play song: " + songs[0])
				monitor.playMusicFile(songs[0])
			}
		})
	}

	robot := gobot.NewRobot("detector",
		[]gobot.Connection{},
		[]gobot.Device{keys},
		work,
	)

	robot.Start()
}

func (monitor *MonitorCmd) Run() {

	monitor.setup()

	//	monitor.waitForExit()
}
