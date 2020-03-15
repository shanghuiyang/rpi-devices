/*
Package dev ...

L298N is an motor driver
which can be used to control the direction and speed of DC motors.

Spec:
           _________________________________________
          |                                         |
          |                                         |
    OUT1 -|                 L298N                   |- OUT3
    OUT2 -|                                         |- OUT4
          |                                         |
          |_________________________________________|
              |   |   |     |   |   |   |   |   |
             12v GND  5V   EN1 IN1 IN2 IN3 IN4 EN2

Pins:
 - OUT1: dc motor A+
 - OUT2: dc motor A-
 - OUT3: dc motor B+
 - OUT4: dc motor B-

 - IN1: input 1 for motor A
 - IN2: input 2 for motor A
 - IN3: input 3 for motor B
 - IN4: input 1 for motor B
 - EN1: enable pin for motor A
 - EN2: enable pin for motor B

*/
package dev

import (
	"github.com/stianeikeland/go-rpio"
)

// L298N ...
type L298N struct {
	in1 rpio.Pin
	in2 rpio.Pin
	in3 rpio.Pin
	in4 rpio.Pin
	ena rpio.Pin
	enb rpio.Pin
}

// NewL298N ...
func NewL298N(in1, in2, in3, in4, ena, enb uint8) *L298N {
	l := &L298N{
		in1: rpio.Pin(in1),
		in2: rpio.Pin(in2),
		in3: rpio.Pin(in3),
		in4: rpio.Pin(in4),
		ena: rpio.Pin(ena),
		enb: rpio.Pin(enb),
	}
	l.in1.Output()
	l.in2.Output()
	l.in3.Output()
	l.in4.Output()
	l.in1.Low()
	l.in2.Low()
	l.in3.Low()
	l.in4.Low()
	l.ena.Pwm()
	l.enb.Pwm()
	l.ena.Freq(64000)
	l.enb.Freq(64000)
	l.Speed(25)
	return l
}

// Forward ...
func (l *L298N) Forward() {
	l.in1.High()
	l.in2.Low()
	l.in3.High()
	l.in4.Low()
}

// Backward ...
func (l *L298N) Backward() {
	l.in1.Low()
	l.in2.High()
	l.in3.Low()
	l.in4.High()
}

// Left ...
func (l *L298N) Left() {
	l.in1.Low()
	l.in2.High()
	l.in3.High()
	l.in4.Low()
}

// Right ...
func (l *L298N) Right() {
	l.in1.High()
	l.in2.Low()
	l.in3.Low()
	l.in4.High()
}

// Stop ...
func (l *L298N) Stop() {
	l.in1.Low()
	l.in2.Low()
	l.in3.Low()
	l.in4.Low()
}

// Speed ...
func (l *L298N) Speed(s uint32) {
	l.ena.DutyCycle(s, 100)
	l.enb.DutyCycle(s, 100)
}
