package main

import (
	"log"
	"time"

	"github.com/shanghuiyang/rpi-devices/dev"
)

const (
	pin = 2
)

func main() {
	sw := dev.NewSW420(pin)
	for {
		shaked := sw.Detected()
		if shaked {
			log.Printf("shaked")
			time.Sleep(1 * time.Second)
			continue
		}
		time.Sleep(100 * time.Millisecond)
	}
}
