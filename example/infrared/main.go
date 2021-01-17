package main

import (
	"log"
	"time"

	"github.com/shanghuiyang/rpi-devices/dev"
	"github.com/shanghuiyang/rpi-devices/util"
	"github.com/stianeikeland/go-rpio"
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

	infr := dev.NewInfrared(pin)
	util.WaitQuit(func() {
		rpio.Close()
	})
	for {
		detectedObj := infr.Detected()
		if detectedObj {
			log.Printf("detected an object")
		} else {
			log.Printf("didn't detect any objects")
		}
		time.Sleep(1 * time.Second)
	}
}
