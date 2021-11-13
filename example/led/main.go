package main

import (
	"time"

	"github.com/shanghuiyang/rpi-devices/dev"
)

const (
	p12 = 26 // led
)

func main() {

	led := dev.NewLedImp(p12)

	for {
		led.On()
		time.Sleep(500 * time.Millisecond)
		led.Off()
		time.Sleep(500 * time.Millisecond)
	}
}
