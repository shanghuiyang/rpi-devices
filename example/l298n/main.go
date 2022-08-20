package main

import (
	"time"

	"github.com/shanghuiyang/rpi-devices/dev"
)

const (
	in1 = 17
	in2 = 23
	in3 = 24
	in4 = 22
	ena = 18
	enb = 13
)

func main() {
	l298n := dev.NewL298N(in1, in2, in3, in4, ena, enb)
	motorA := dev.NewDCMotor(l298n.MotorA)
	motorB := dev.NewDCMotor(l298n.MotorB)

	motorA.SetSpeed(30)
	motorB.SetSpeed(60)

	motorA.Forward()
	motorB.Backward()
	time.Sleep(5 * time.Second)
	motorA.Stop()
	motorB.Stop()
}
