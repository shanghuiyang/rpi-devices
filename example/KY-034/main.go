package main

import (
	"github.com/shanghuiyang/rpi-devices/dev"
	"github.com/stianeikeland/go-rpio"
	"log"
	"time"
)

func main(){
	if err := rpio.Open(); err != nil {
		log.Fatalf("failed to open rpio, error: %v", err)
		return
	}
	defer rpio.Close()

	f := dev.NewSevenColorLed(18)

	for {
		f.On()
		time.Sleep(time.Second)
		f.Off()
		time.Sleep(time.Second)
	}
}
