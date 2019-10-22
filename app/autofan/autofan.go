package main

import (
	"log"
	"time"

	"github.com/shanghuiyang/rpi-devices/dev"
	"github.com/stianeikeland/go-rpio"
)

const (
	relayPin           = 7
	intervalTime       = 1 * time.Minute
	triggerTemperature = 27.3
)

func main() {
	if err := rpio.Open(); err != nil {
		log.Fatalf("failed to open rpio, error: %v", err)
		return
	}
	defer rpio.Close()

	t := dev.NewTemperature()
	if t == nil {
		log.Printf("failed to new a temperature device")
		return
	}

	r := dev.NewRelay(relayPin)
	if r == nil {
		log.Printf("failed to new a relay device")
		return
	}

	f := &autoFan{
		temperature: t,
		relay:       r,
	}
	f.start()
}

type autoFan struct {
	temperature *dev.Temperature
	relay       *dev.Relay
}

func (f *autoFan) start() {
	for {
		time.Sleep(intervalTime)
		c, err := f.temperature.GetTemperature()
		if err != nil {
			log.Printf("failed to get temperature, error: %v", err)
			continue
		}
		if c >= triggerTemperature {
			f.on()
		} else {
			f.off()
		}
	}
}

func (f *autoFan) on() {
	f.relay.On()
}

func (f *autoFan) off() {
	f.relay.Off()
}
