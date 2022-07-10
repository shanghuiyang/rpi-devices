package main

import (
	"log"
	"time"

	"github.com/shanghuiyang/rpi-devices/dev"
)

const (
	pin = 26
)

func main() {

	f := dev.NewFan(pin)

	for {
		f.On()
		log.Print("on")
		time.Sleep(5 * time.Second)
		f.Off()
		log.Print("off")
		time.Sleep(5 * time.Second)
	}
}
