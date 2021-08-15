package main

import (
	"fmt"
	"log"
	"time"

	"github.com/shanghuiyang/rpi-devices/dev"
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
		d, err := hcsr04.Dist()
		if err != nil {
			log.Printf("failed to get distance, error: %v", err)
			time.Sleep(100 * time.Millisecond)
			continue
		}
		fmt.Printf("%.2f cm\n", d)
		time.Sleep(1 * time.Second)
	}
}
