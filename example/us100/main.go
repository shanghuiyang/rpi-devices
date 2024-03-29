package main

import (
	"fmt"
	"log"
	"time"

	"github.com/shanghuiyang/rpi-devices/dev"
)

const (
	trig = 21
	echo = 26

	us100Dev = "/dev/ttyAMA0"
	baud     = 9600
)

func main() {
	var mode int
	fmt.Printf("choose interface type: 1: GPIO, 2: UART\n")
	fmt.Printf("mode: ")
	if n, err := fmt.Scanf("%d", &mode); n != 1 || err != nil {
		log.Printf("invalid operator, error: %v", err)
		return
	}
	if mode != 1 && mode != 2 {
		log.Printf("invalid input")
		return
	}

	// gpio interface
	if mode == 1 {
		u, err := dev.NewUS100GPIO(trig, echo)
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

	// uart interface
	u, err := dev.NewUS100UART(us100Dev, baud)
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
