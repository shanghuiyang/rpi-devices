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
	ads, err := NewADS1015()
	if err != nil {
		return nil, err
	}
	j := &Joystick{
		swPin: rpio.Pin(sw),
		ads:   ads,
	}
	j.swPin.Input()
	return j, nil
}

// X ...
func (j *Joystick) X() (x float64) {
	v, err := j.ads.Read(0)
	if err != nil {
		return 0
	}
	return v
}

// Y ...
func (j *Joystick) Y() (y float64) {
	v, err := j.ads.Read(1)
	if err != nil {
		return 0
	}
	return v
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
