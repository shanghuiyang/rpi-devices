package main

import (
	"fmt"
	"log"

	"github.com/shanghuiyang/rpi-devices/devices"
)

const (
	p8  = 8  // in1 for step motor
	p25 = 25 // in2 for step motor
	p24 = 24 // in3 for step motor
	p23 = 23 // in4 for step motor
)

func main() {
	m := devices.NewStepMotor(p8, p25, p24, p23)
	go m.Start()
	log.Printf("step motor is ready for service\n")

	var input int
	for {
		fmt.Printf(">>op: ")
		if n, err := fmt.Scanf("%d", &input); n != 1 || err != nil {
			log.Printf("invalid operator, error: %v", err)
			continue
		}
		if input == 0 {
			break
		}
		op := devices.Operator(input)
		devices.ChStepMotorOp <- op
	}
	log.Printf("step motor stop service\n")
}
