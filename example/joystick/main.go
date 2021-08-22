package main

import (
	"log"
	"time"

	"github.com/shanghuiyang/rpi-devices/dev"
)

const (
	swPin = 7
)

func main() {
	j, err := dev.NewJoystickImp(swPin)
	if err != nil {
		log.Printf("failed to new joystick")
		return
	}

	for {
		x := j.X()
		y := j.Y()
		z := j.Z()
		log.Printf("x: %v, y: %v, z: %v", x, y, z)
		time.Sleep(300 * time.Millisecond)
	}
}
