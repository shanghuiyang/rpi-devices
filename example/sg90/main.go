package main

import (
	"fmt"
	"log"

	"github.com/shanghuiyang/rpi-devices/dev"
)

const (
	pin = 18
)

func main() {
	sg := dev.NewSG90(pin)
	var angle float64
	for {
		fmt.Printf(">>angle: ")
		if n, err := fmt.Scanf("%f", &angle); n != 1 || err != nil {
			log.Printf("invalid angle, error: %v", err)
			continue
		}
		sg.Roll(angle)
	}
}
