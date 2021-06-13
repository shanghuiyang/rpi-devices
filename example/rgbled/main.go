package main

import (
	"log"
	"time"

	"github.com/jakefau/rpi-devices/dev"
	"github.com/stianeikeland/go-rpio"
)

const (
	r = 10 // red
	g = 9  // green
	b = 11 // blue
)

func main() {
	if err := rpio.Open(); err != nil {
		log.Fatalf("failed to open rpio, error: %v", err)
		return
	}
	defer rpio.Close()

	rgbled := dev.NewRGBLed(r, g, b)

	rgbled.RedOn()
	time.Sleep(5 * time.Second)
	rgbled.RedOff()
	rgbled.GreenOn()
	time.Sleep(5 * time.Second)
	rgbled.GreenOff()
	rgbled.BlueOn()
	time.Sleep(5 * time.Second)
	rgbled.RedOn()
	rgbled.GreenOn()
	time.Sleep(5 * time.Second)
	rgbled.RedOff()
	rgbled.GreenOff()
	rgbled.BlueOff()

}
