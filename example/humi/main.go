package main

import (
	"log"
	"time"

	"github.com/shanghuiyang/rpi-devices/dev"
	"github.com/shanghuiyang/rpi-devices/util"
	"github.com/stianeikeland/go-rpio"
)

const (
	pin = 17
)

func main() {
	if err := rpio.Open(); err != nil {
		log.Fatalf("failed to open rpio, error: %v", err)
		return
	}
	defer rpio.Close()

	h := dev.NewHumi(pin)
	util.WaitQuit(func() {
		rpio.Close()
	})
	for {
		if h.Detected() {
			log.Printf("detected humidity")
		}
		time.Sleep(1 * time.Second)
	}
}
