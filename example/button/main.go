package main

import (
	"log"
	"time"

	"github.com/shanghuiyang/rpi-devices/dev"
	"github.com/shanghuiyang/rpi-devices/util"
	"github.com/stianeikeland/go-rpio"
)

const (
	pin    = 17
	pinLed = 26
)

func main() {
	if err := rpio.Open(); err != nil {
		log.Fatalf("failed to open rpio, error: %v", err)
		return
	}
	defer rpio.Close()

	led := dev.NewLedImp(pinLed)
	btn := dev.NewButtonImp(pin)
	util.WaitQuit(func() {
		rpio.Close()
	})

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
