package main

import (
	"fmt"

	"github.com/shanghuiyang/rpi-devices/dev"
)

const (
	in1 = 12
	in2 = 16
	in3 = 20
	in4 = 21
)

func main() {
	var step int
	var angle float64
	motor := dev.NewBYJ2848(in1, in2, in3, in4)

	for {
		fmt.Printf(">>steps: ")
		if n, err := fmt.Scanf("%d", &step); n != 1 || err != nil {
			fmt.Printf("invalid steps, error: %v", err)
			continue
		}
		if step == 0 {
			break
		}
		motor.Step(step)

		fmt.Printf(">>angle: ")
		if n, err := fmt.Scanf("%f", &angle); n != 1 || err != nil {
			fmt.Printf("invalid angle, error: %v", err)
			continue
		}
		if angle == 0 {
			break
		}
		motor.Roll(angle)
	}
}
