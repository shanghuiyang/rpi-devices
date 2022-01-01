package main

import (
	"fmt"
	"log"
	"time"

	"github.com/shanghuiyang/rpi-devices/dev"
)

const (
	in1 = 12
	in2 = 16
	in3 = 20
	in4 = 21

	steps  = 30
	feedAt = "10:00"
)

func main() {
	var h, m int
	if n, err := fmt.Sscanf(feedAt, "%d:%d", &h, &m); n != 2 || err != nil {
		log.Fatalf("parse feed time error: %v", err)
	}

	motor := dev.NewBYJ2848(in1, in2, in3, in4)
	for {
		now := time.Now()
		if now.Hour() == h && now.Minute() == m {
			motor.Step(steps)
			log.Print("fed")
		}
		time.Sleep(time.Minute)
	}
}
