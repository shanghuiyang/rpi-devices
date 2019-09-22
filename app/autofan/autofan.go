package main

import (
	"log"
	"time"

	dev "github.com/shanghuiyang/rpi-devices/devices"
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

	f := &fan{
		temperature: t,
		relay:       r,
	}
	f.start()
}

type fan struct {
	temperature *dev.Temperature
	relay       *dev.Relay
}

func (f *fan) start() {
	for {
		time.Sleep(intervalTime)
		c, err := f.temperature.GetTemperature()
		if err != nil {
			log.Printf("failed to get temperature, error: %v", err)
			continue
		}
		if c >= triggerTemperature {
			f.relay.On()
		} else {
			f.relay.Off()
		}
	}
}
