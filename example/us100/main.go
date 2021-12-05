package main

import (
	"fmt"
	"log"
	"time"

	"github.com/shanghuiyang/rpi-devices/dev"
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
		u, err := dev.NewUS100(&dev.US100Config{
			Mode: dev.TTLMode,
			Trig: 21,
			Echo: 26,
		})
		if err != nil {
			log.Fatalf("new us100 error: %v", err)
		}
		for {
			dist, err := u.Dist()
			if err != nil {
				log.Printf("failed to get distance")
				continue
			}
			log.Printf("%.2f cm\n", dist)
			time.Sleep(50 * time.Millisecond)
		}
	}

	// uart mode
	u, err := dev.NewUS100(&dev.US100Config{
		Mode: dev.UartMode,
		Dev:  "/dev/ttyAMA0",
		Baud: 9600,
	})
	if err != nil {
		log.Fatalf("new us100 error: %v", err)
	}
	defer u.Close()

	for {
		dist, err := u.Dist()
		if err != nil {
			log.Printf("failed to get distance")
			continue
		}
		log.Printf("%.2f cm\n", dist)
		time.Sleep(50 * time.Millisecond)
	}
}
