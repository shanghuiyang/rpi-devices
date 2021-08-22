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
	c := dev.NewCollisionDetector(pin)
	for {
		collided := c.Detected()
		if collided {
			log.Printf("collided")
		}
		time.Sleep(10 * time.Millisecond)
	}
}
