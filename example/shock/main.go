package main

import (
	"github.com/shanghuiyang/rpi-devices/dev"
	"github.com/shanghuiyang/rpi-devices/util"
	"github.com/stianeikeland/go-rpio"
	"log"
	"time"
)

const (
	pin = 18
)

func main() {
	if err := rpio.Open(); err != nil {
		log.Fatalf("failed to open rpio, error: %v", err)
		return
	}
	defer rpio.Close()
	s := dev.NewShockSensor(pin)
	util.WaitQuit(func() {
		rpio.Close()
	})
	for {
		shocked := s.Shock()
		if shocked {
			log.Printf("shocked")
		}
		time.Sleep(10 * time.Millisecond)
	}

}