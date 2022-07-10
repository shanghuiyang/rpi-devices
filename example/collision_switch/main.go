package main

import (
	"log"
	"time"

	"github.com/shanghuiyang/rpi-devices/dev"
)

const (
	pin = 12
)

func main() {
	c := dev.NewCollisionSwitch(pin)
	for {
		if !c.Detected() {
			time.Sleep(100 * time.Millisecond)
			continue
		}
		log.Printf("collided")
		time.Sleep(1 * time.Second)
	}
}
