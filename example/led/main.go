package main

import (
	"time"

	"github.com/shanghuiyang/rpi-devices/dev"
)

const (
	pin = 26
)

func main() {

	led := dev.NewLedImp(pin)

	for {
		led.On()
		time.Sleep(500 * time.Millisecond)
		led.Off()
		time.Sleep(500 * time.Millisecond)
	}
}
