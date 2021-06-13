package main

import (
	"log"
	"time"

	"github.com/jakefau/rpi-devices/dev"
	"github.com/jakefau/rpi-devices/util"
	"github.com/stianeikeland/go-rpio"
)

const (
	pin = 12
)

func main() {
	if err := rpio.Open(); err != nil {
		log.Fatalf("failed to open rpio, error: %v", err)
		return
	}
	defer rpio.Close()

	c := dev.NewCollision(pin)
	util.WaitQuit(func() {
		rpio.Close()
	})
	for {
		collided := c.Collided()
		if collided {
			log.Printf("collided")
		}
		time.Sleep(10 * time.Millisecond)
	}
}
