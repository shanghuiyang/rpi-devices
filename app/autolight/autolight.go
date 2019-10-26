package main

import (
	"log"
	"time"

	"github.com/shanghuiyang/rpi-devices/dev"
	"github.com/stianeikeland/go-rpio"
)

const (
	voicePin = 2
	ledPin   = 16
)

func main() {
	if err := rpio.Open(); err != nil {
		log.Fatalf("failed to open rpio, error: %v", err)
		return
	}
	defer rpio.Close()

	v := dev.NewVoiceDetector(voicePin)
	if v == nil {
		log.Printf("failed to new a voice detector")
		return
	}

	light := dev.NewLed(ledPin)
	if light == nil {
		log.Printf("failed to new a led light")
		return
	}

	l := &autoLight{
		voice: v,
		light: light,
	}
	l.start()
}

type autoLight struct {
	voice *dev.VoiceDetector
	light *dev.Led
}

func (a *autoLight) start() {
	lastTrig := time.Now()
	isLightOn := false
	for {
		time.Sleep(200 * time.Millisecond)
		if a.voice.Detected() {
			log.Printf("detected a voice")
			if !isLightOn {
				a.light.On()
				isLightOn = true
			}
			lastTrig = time.Now()
			continue
		}
		if time.Now().Sub(lastTrig).Seconds() > 35 && isLightOn {
			log.Printf("timeout, light off")
			a.light.Off()
			isLightOn = false
		}
	}
}
