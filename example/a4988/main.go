package main

import (
	"fmt"

	"github.com/shanghuiyang/rpi-devices/dev"
)

const (
	step = 26
	dir  = 12
	ms1  = 22
	ms2  = 27
	ms3  = 17
	enb  = 25
)

func main() {
	stepper := dev.NewA4988(step, dir, ms1, ms2, ms3)
	stepper.SetMode(dev.HalfMode)
	var s int
	for {
		fmt.Printf(">>steps: ")
		if n, err := fmt.Scanf("%d", &s); n != 1 || err != nil {
			fmt.Printf("invalid steps, error: %v", err)
			continue
		}
		if s == 0 {
			break
		}
		stepper.Step(s)
	}
}
