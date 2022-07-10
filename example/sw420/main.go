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
		if !sw.Detected() {
			time.Sleep(100 * time.Millisecond)
			continue
		}
		log.Printf("shaking")
		time.Sleep(1 * time.Second)
	}
}
