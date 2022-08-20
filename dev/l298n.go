/*
L298N is a motor driver used to control the direction and speed of DC motors.

Spec:
           _________________________________________
          |                                         |
          |                                         |
    OUT1 -|                 L298N                   |- OUT3
    OUT2 -|                                         |- OUT4
          |                                         |
          |_________________________________________|
              |   |   |     |   |   |   |   |   |
             12v GND  5V   ENA IN1 IN2 IN3 IN4 ENB

Pins:
 - OUT1: dc motor A+
 - OUT2: dc motor A-
 - OUT3: dc motor B+
 - OUT4: dc motor B-

 - 12v: +battery
 - GND: -battery (and any gnd pin of raspberry pi if motors and raspberry pi use different battery sources)
 - IN1: any data pin
 - IN2: any data pin
 - IN3: any data pin
 - IN4: any data pin
 - EN1: must be one of GPIO 12, 13, 18 or 19 (pwn pins)
 - EN2: must be one of GPIO 12, 13, 18 or 19 (pwn pins)

*/
package dev

import (
	"github.com/stianeikeland/go-rpio/v4"
)

// L298N implements MotorDriver interface
type L298N struct {
	MotorA MotorDriver
	MotorB MotorDriver
}

type l298nMotorDriver struct {
	in1 rpio.Pin
	in2 rpio.Pin
	en  rpio.Pin
}

// NewL298N ...
func NewL298N(in1, in2, in3, in4, ena, enb uint8) *L298N {
	l := &L298N{
		MotorA: newL298NMotorDriver(in1, in2, ena),
		MotorB: newL298NMotorDriver(in3, in4, enb),
	}
	return l
}

func newL298NMotorDriver(in1, in2, en uint8) *l298nMotorDriver {
	m := &l298nMotorDriver{
		in1: rpio.Pin(in1),
		in2: rpio.Pin(in2),
		en:  rpio.Pin(en),
	}
	m.in1.Output()
	m.in2.Output()
	m.in1.Low()
	m.in2.Low()
	m.en.Pwm()
	m.en.Freq(50 * 100)
	return m
}

// Forward ...
func (m *l298nMotorDriver) Forward() {
	m.in1.High()
	m.in2.Low()
}

// Backward ...
func (m *l298nMotorDriver) Backward() {
	m.in1.Low()
	m.in2.High()
}

// Stop ...
func (m *l298nMotorDriver) Stop() {
	m.in1.Low()
	m.in2.Low()
}

// Speed ...
func (m *l298nMotorDriver) SetSpeed(s uint32) {
	m.en.DutyCycle(0, 100)
	m.en.DutyCycle(s, 100)
}
