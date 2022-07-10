package main

import (
	"log"
	"time"

	"github.com/shanghuiyang/rpi-devices/dev"
)

const (
	pin = 21
)

func main() {
	btn := dev.NewButtonImp(pin)
	for {
		if !btn.Pressed() {
			time.Sleep(100 * time.Millisecond)
			continue
		}
		log.Printf("the button was pressed")
		time.Sleep(1 * time.Second)
	}
}
