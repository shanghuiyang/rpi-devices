package main

import (
	"log"
	"time"

	"github.com/shanghuiyang/rpi-devices/dev"
)

const (
	pin = 17
)

func main() {
	h := dev.NewHumidityDetector(pin)
	for {
		if h.Detected() {
			log.Printf("detected humidity")
		}
		time.Sleep(1 * time.Second)
	}
}
