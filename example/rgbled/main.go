package main

import (
	"github.com/shanghuiyang/rpi-devices/dev"
	"github.com/stianeikeland/go-rpio"
	"log"
	"time"
)

const (
	r = 14 // red
	g = 15 // green
	b = 18 // blue
)

func main(){
	if err := rpio.Open(); err != nil {
		log.Fatalf("failed to open rpio, error: %v", err)
		return
	}
	defer rpio.Close()

	rgbled := dev.NewRGBLed(r,g,b)

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