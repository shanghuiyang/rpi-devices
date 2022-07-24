package main

import (
	"time"

	"github.com/shanghuiyang/rpi-devices/dev"
)

const pin = 26

func main() {
	led := dev.NewLedImp(pin)
	for {
		led.On()
		time.Sleep(1 * time.Second)
		led.Off()
		time.Sleep(1 * time.Second)
	}
}
