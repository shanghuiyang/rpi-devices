package main

import (
	"log"
	"time"

	"github.com/shanghuiyang/rpi-devices/dev"
)

const (
	pin    = 21
	pinLed = 26
)

func main() {
	led := dev.NewLedImp(pinLed)
	btn := dev.NewButtonImp(pin)

	on := false
	for {
		pressed := btn.Pressed()
		if pressed {
			log.Printf("the button was pressed")
			if on {
				on = false
				led.Off()
			} else {
				led.On()
				on = true
			}
		}
		time.Sleep(300 * time.Millisecond)
	}
}
