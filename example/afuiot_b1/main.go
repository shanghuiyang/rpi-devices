package main

import (
	"log"
	"time"

	"github.com/shanghuiyang/rpi-devices/dev"
)

const (
	devName = "/dev/ttyAMA0"
	baud    = 38400
)

func main() {
	b1, err := dev.NewAfuiotB1(devName, baud)
	if err != nil {
		log.Printf("failed to create irt, error: %v", err)
		return
	}
	defer b1.Close()

	if err := b1.Set(dev.SetAfuiotB1HumanTemp); err != nil {
		log.Printf("set failed: %v", err)
	}

	for {
		time.Sleep(100 * time.Millisecond)
		t, err := b1.Temperature()
		if err != nil {
			log.Printf("failed to get temperature, error: %v", err)
			continue
		}
		log.Printf("%.1f", t)
	}
}
