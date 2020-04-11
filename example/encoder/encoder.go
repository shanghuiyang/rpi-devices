package main

import (
	"log"

	"github.com/shanghuiyang/rpi-devices/dev"
	"github.com/stianeikeland/go-rpio"
)

const (
	pinEncoder = 6
)

func main() {
	if err := rpio.Open(); err != nil {
		log.Fatalf("failed to open rpio, error: %v", err)
		return
	}
	defer rpio.Close()

	e := dev.NewEncoder(pinEncoder)
	count := 0
	for {
		count += e.Count1()
		log.Printf("%v", count)
	}
}
