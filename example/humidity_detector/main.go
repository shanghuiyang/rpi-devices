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
		if !h.Detected() {
			time.Sleep(100 * time.Millisecond)
			continue
		}
		log.Printf("detected humidity")
		time.Sleep(5 * time.Second)
	}
}
