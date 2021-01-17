package main

import (
	"fmt"
	"log"
	"time"

	"github.com/shanghuiyang/rpi-devices/dev"
	"github.com/stianeikeland/go-rpio"
)

func main() {
	var mode int
	fmt.Printf("Please choose mode: 1: TTL, 2: Uart\n")
	fmt.Printf("mode: ")
	if n, err := fmt.Scanf("%d", &mode); n != 1 || err != nil {
		log.Printf("invalid operator, error: %v", err)
		return
	}
	if mode != 1 && mode != 2 {
		log.Printf("invalid input")
		return
	}

	// ttl mode
	if mode == 1 {
		if err := rpio.Open(); err != nil {
			log.Fatalf("failed to open rpio, error: %v", err)
			return
		}
		defer rpio.Close()

		u := dev.NewUS100(&dev.US100Config{
			Mode: dev.TTLMode,
			Trig: 21,
			Echo: 26,
		})
		for {
			dist := u.Dist()
			if dist < 0 {
				log.Printf("failed to get distance")
				continue
			}
			log.Printf("%.2f cm\n", dist)
			time.Sleep(50 * time.Millisecond)
		}
	}

	// uart mode
	u := dev.NewUS100(&dev.US100Config{
		Mode: dev.UartMode,
		Dev:  "/dev/ttyAMA0",
		Baud: 9600,
	})
	defer u.Close()

	for {
		dist := u.Dist()
		if dist < 0 {
			log.Printf("failed to get distance")
			continue
		}
		log.Printf("%.2f cm\n", dist)
		time.Sleep(50 * time.Millisecond)
	}
}
