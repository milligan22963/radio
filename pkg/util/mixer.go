// Package util contains utilities
package util

import (
	"fmt"
	"os"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/wav"
	"github.com/sirupsen/logrus"
)

// RadioMixer is a struct defining a mixing radio stream
type RadioMixer struct {
	staticStreamers   []beep.StreamSeekCloser
	staticStreamIndex int
	stations          map[float64]beep.StreamSeekCloser
	currentStation    beep.StreamSeekCloser
}

func (rm *RadioMixer) AddStatic(filename string) error {
	musicFile, err := os.Open(filename)
	if err != nil {
		return err
	}

	streamer, format, err := wav.Decode(musicFile)
	if err != nil {
		return err
	}

	if format.SampleRate.N(time.Second) != 44100 {
		return fmt.Errorf("unsupported sample rate")
	}

	rm.staticStreamers = append(rm.staticStreamers, streamer)
	return nil
}

func (rm *RadioMixer) AddStation(channel float64, filename string) error {
	musicFile, err := os.Open(filename)
	if err != nil {
		return err
	}

	streamer, format, err := mp3.Decode(musicFile)
	if err != nil {
		return err
	}

	if format.SampleRate.N(time.Second) != 44100 {
		return fmt.Errorf("unsupported sample rate")
	}

	rm.stations[channel] = streamer
	return nil
}

func (rm *RadioMixer) PlayStation(channel float64) {
	//  If this station has an entry, then play it
	if val, ok := rm.stations[channel]; ok {
		rm.currentStation = val
	} else {
		rm.currentStation = nil
	}
}

func (rm *RadioMixer) Stream(samples [][2]float64) (n int, ok bool) {
	// return up to sampleCount, if we have less then pad w/ zeros
	// or loop
	sampleCount := len(samples)
	samplesReturned := 0
	for samplesReturned < sampleCount {
		if rm.currentStation != nil {
			logrus.Info("playing current station")
			// stream station
			n, ok = rm.currentStation.Stream(samples[samplesReturned:])
			// If it's drained, we start over
			if !ok {
				rm.currentStation.Seek(0)
			}
		} else {
			// We stream from the current stream in the array of static
			n, ok = rm.staticStreamers[rm.staticStreamIndex].Stream(samples[samplesReturned:])
			// If it's drained, we move to the next one or wrap to start over
			if !ok {
				rm.staticStreamIndex += 1
				if rm.staticStreamIndex >= len(rm.staticStreamers) {
					rm.staticStreamIndex = 0
				}
				// start at the beginning
				err := rm.staticStreamers[rm.staticStreamIndex].Seek(0)
				if err != nil {
					logrus.Errorf("seeking error: +%v", err)
				}
			}
		}
		// We update the number of returned samples.
		samplesReturned += n
	}
	return sampleCount, true
}

// Err generates a new error for now
func (rm *RadioMixer) Err() error {
	return fmt.Errorf("failed playing mixer")
}

func (rm *RadioMixer) Close() {
	for _, v := range rm.staticStreamers {
		v.Close()
	}

	for _, v := range rm.stations {
		v.Close()
	}
}

func (rm *RadioMixer) Initialize() {
	rm.staticStreamIndex = 0
	rm.stations = make(map[float64]beep.StreamSeekCloser)
}
