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
		fmt.Printf("failed to open file: %v\n", err)
	}

	streamer, format, err := mp3.Decode(f)
	if err != nil {
		fmt.Printf("mp3 decoding fail: %v\n", err)
	}
	defer streamer.Close()
	err = speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	if err != nil {
		fmt.Printf("speaker init fail: %v\n", err)
	}
	done := make(chan bool)
	speaker.Play(beep.Seq(streamer, beep.Callback(func() {
		fmt.Println("song has ended.")
		done <- true
	})))

	<-done
}

func (monitor *MonitorCmd) setup() {

	// configure gpios
	keys := keyboard.NewDriver()

	work := func() {
		keys.On(keyboard.Key, func(data interface{}) {
			key := data.(keyboard.KeyEvent)

			if key.Key == keyboard.A {
				fmt.Println("starting to play song...")
				monitor.playMusicFile("Abbott and Costello Whos On First.mp3")
			} else {
				fmt.Println("keyboard event!", key, key.Char)
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
