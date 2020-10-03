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
	pcf   *PCF8591
}

// NewJoystick ...
func NewJoystick(sw uint8) (*Joystick, error) {
	p, err := NewPCF8591()
	if err != nil {
		return nil, err
	}
	return &Joystick{
		swPin: rpio.Pin(sw),
		pcf:   p,
	}, nil
}

// X ...
// -123 =< x <= 132
// x > 0: left
// x = 0: home
// x < 0: right
func (j *Joystick) X() (x int) {
	data := j.pcf.ReadAIN0()
	if len(data) == 0 {
		return 0
	}
	x = int(data[0]) - 123
	return
}

// Y ...
// -124 =< y <= 131
// y > 0: up
// y = 0: home
// y < 0: down
func (j *Joystick) Y() (y int) {
	data := j.pcf.ReadAIN1()
	if len(data) == 0 {
		return 0
	}
	y = 131 - int(data[0])
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
