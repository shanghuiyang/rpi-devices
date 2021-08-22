package main

import (
	"fmt"
	"log"

	"github.com/shanghuiyang/rpi-devices/dev"
)

const (
	p8  = 8  // in1 for step motor
	p25 = 25 // in2 for step motor
	p24 = 24 // in3 for step motor
	p23 = 23 // in4 for step motor
)

func main() {
	var angle float64
	motor := dev.NewBYJ2848(p8, p25, p24, p23)
	for {
		fmt.Printf(">>op: ")
		if n, err := fmt.Scanf("%f", &angle); n != 1 || err != nil {
			log.Printf("invalid angle, error: %v", err)
			continue
		}
		if angle == 0 {
			break
		}
		motor.Roll(angle)
	}
	log.Printf("step motor stop service\n")
}
