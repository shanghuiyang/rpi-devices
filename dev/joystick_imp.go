/*
JoystickImp is 2-axis joystick module.

Connect to Pi:
 - +V5: any 5v
 - GND: any gnd pin
 - SM : any data pin
 - Rx: ADS1015->A0
 - Ry: ADS1015->A1
*/
package dev

import (
	"github.com/stianeikeland/go-rpio/v4"
)

// JoystickImp ...
type JoystickImp struct {
	swPin rpio.Pin
	ads   *ADS1015
}

// NewJoystickImp ...
func NewJoystickImp(sw uint8) (*JoystickImp, error) {
	ads, err := NewADS1015()
	if err != nil {
		return nil, err
	}
	j := &JoystickImp{
		swPin: rpio.Pin(sw),
		ads:   ads,
	}
	j.swPin.Input()
	return j, nil
}

// X ...
func (j *JoystickImp) X() float64 {
	x, err := j.ads.Read(0)
	if err != nil {
		return 0
	}
	return x
}

// Y ...
func (j *JoystickImp) Y() float64 {
	y, err := j.ads.Read(1)
	if err != nil {
		return 0
	}
	return y
}

// Z ...
// z = 1: pressed
// z = 0: home
func (j *JoystickImp) Z() int {
	if j.swPin.Read() == rpio.Low {
		return 1 // pressed
	}
	return 0 // home
}
