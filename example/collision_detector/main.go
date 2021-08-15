package main

import (
	"log"
	"time"

	"github.com/shanghuiyang/rpi-devices/dev"
	"github.com/shanghuiyang/rpi-devices/util"
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

	c := dev.NewCollisionDetector(pin)
	util.WaitQuit(func() {
		rpio.Close()
	})
	for {
		collided := c.Detected()
		if collided {
			log.Printf("collided")
		}
		time.Sleep(10 * time.Millisecond)
	}
}
