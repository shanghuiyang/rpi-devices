package main

import (
	"log"
	"time"

	"github.com/shanghuiyang/rpi-devices/base"
	"github.com/shanghuiyang/rpi-devices/dev"
	"github.com/stianeikeland/go-rpio"
)

const (
	swPin = 4
)

func main() {
	if err := rpio.Open(); err != nil {
		log.Fatalf("failed to open rpio, error: %v", err)
		return
	}
	defer rpio.Close()

	j, err := dev.NewJoystick(swPin)
	if err != nil {
		log.Printf("failed to new joystick")
		return
	}
	base.WaitQuit(func() {
		rpio.Close()
	})

	for {
		x := j.X()
		y := j.Y()
		z := j.Z()
		log.Printf("x: %v, y: %v, z: %v", x, y, z)
		time.Sleep(100 * time.Millisecond)
	}
}
