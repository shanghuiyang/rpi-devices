package main

import (
	"log"
	"time"

	"github.com/shanghuiyang/rpi-devices/base"
	"github.com/shanghuiyang/rpi-devices/dev"
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

	cswitch := dev.NewCollisionSwitch(pin)
	base.WaitQuit(func() {
		rpio.Close()
	})
	for {
		collided := cswitch.Collided()
		if collided {
			log.Printf("collided")
		}
		time.Sleep(10 * time.Millisecond)
	}
}
