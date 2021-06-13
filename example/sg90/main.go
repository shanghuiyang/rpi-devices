package main

import (
	"fmt"
	"log"

	"github.com/jakefau/rpi-devices/dev"
	"github.com/stianeikeland/go-rpio"
)

const (
	p18 = 18
)

func main() {
	if err := rpio.Open(); err != nil {
		log.Fatalf("failed to open rpio, error: %v", err)
		return
	}
	defer rpio.Close()

	sg := dev.NewSG90(p18)
	var angle int
	for {
		fmt.Printf(">>angle: ")
		if n, err := fmt.Scanf("%d", &angle); n != 1 || err != nil {
			log.Printf("invalid angle, error: %v", err)
			continue
		}
		sg.Roll(angle)
	}
	log.Printf("quit")
}
