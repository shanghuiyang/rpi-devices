package main

import (
	"fmt"
	"log"
	"time"

	"github.com/jakefau/rpi-devices/dev"
	"github.com/stianeikeland/go-rpio"
)

const (
	pinTrig = 21
	pinEcho = 26
)

func main() {
	if err := rpio.Open(); err != nil {
		log.Fatalf("failed to open rpio, error: %v", err)
		return
	}
	defer rpio.Close()

	hcsr04 := dev.NewHCSR04(pinTrig, pinEcho)
	for {
		dist := hcsr04.Dist()
		fmt.Printf("%.2f cm\n", dist)
		time.Sleep(1 * time.Second)
	}
}
