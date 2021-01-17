package main

import (
	"log"
	"time"

	"github.com/shanghuiyang/rpi-devices/dev"
	"github.com/shanghuiyang/rpi-devices/util"
	"github.com/stianeikeland/go-rpio"
)

const (
	pin = 2
)

func main() {
	if err := rpio.Open(); err != nil {
		log.Fatalf("failed to open rpio, error: %v", err)
		return
	}
	defer rpio.Close()

	sw := dev.NewSW420(pin)
	util.WaitQuit(func() {
		rpio.Close()
	})

	for {
		shaked := sw.Shaked()
		if shaked {
			log.Printf("shaked")
			time.Sleep(1 * time.Second)
			continue
		}
		time.Sleep(100 * time.Millisecond)
	}
}
