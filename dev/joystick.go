/*
Package dev ...

Connect to Pi:
 - +V5: any 5v
 - GND: any gnd pin
 - SM : any data pin
 - Rx: PCF8591->AIN0
 - Ry: PCF8591->AIN1
*/

package dev

import (
	"github.com/stianeikeland/go-rpio"
)

// Joystick ...
type Joystick struct {
	swPin rpio.Pin
	ads   *ADS1015
}

// NewJoystick ...
func NewJoystick(sw uint8) (*Joystick, error) {
	// p, err := NewPCF8591()
	ads, err := NewADS1015()
	if err != nil {
		return nil, err
	}
	return &Joystick{
		swPin: rpio.Pin(sw),
		ads:   ads,
	}, nil
}

// X ...
// -123 =< x <= 132
// x > 0: left
// x = 0: home
// x < 0: right
func (j *Joystick) X() (x int) {
	v, err := j.ads.Read(0)
	if err != nil {
		return 0
	}
	x = int(v)
	return
}

// Y ...
// -124 =< y <= 131
// y > 0: up
// y = 0: home
// y < 0: down
func (j *Joystick) Y() (y int) {
	v, err := j.ads.Read(1)
	if err != nil {
		return 0
	}
	y = int(v)
	return
}

// Z ...
// z = 1: pressed
// z = 0: home
func (j *Joystick) Z() (z int) {
	if j.swPin.Read() == rpio.Low {
		return 1 // pressed
	}
	return 0 // home
}
