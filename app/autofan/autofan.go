package main

import (
	"log"
	"time"

	dev "github.com/shanghuiyang/rpi-devices/devices"
)

const (
	relayPin           = 7
	intervalTime       = 1 * time.Minute
	triggerTemperature = 27.3
)

func main() {
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
