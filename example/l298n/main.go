package main

import (
	"time"

	"github.com/shanghuiyang/rpi-devices/dev"
)

const (
	in1 = 17
	in2 = 23
	in3 = 27
	in4 = 22
	ena = 13
	enb = 19
)

func main() {
	l298n := dev.NewL298N(in1, in2, in3, in4, ena, enb)
	motorA := dev.NewDCMotor(l298n.MotorA)
	motorB := dev.NewDCMotor(l298n.MotorB)

	motorA.SetSpeed(25)
	motorB.SetSpeed(25)

	motorA.Forward()
	motorB.Forward()
	time.Sleep(2 * time.Second)
	motorA.Stop()
	motorB.Stop()
}
