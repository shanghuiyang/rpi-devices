/*
Package dev ...

Connect to Pi:
 - +5v: any 5v pin
 - gnd: any gnd pin
 - out1: any data pin
 - out2: any data pin
 - out3: any data pin
 - out4: any data pin
 - key: null

*/
package dev

import (
	"github.com/stianeikeland/go-rpio"
)

// RX480E4 ...
type RX480E4 struct {
	out1 rpio.Pin
	out2 rpio.Pin
	out3 rpio.Pin
	out4 rpio.Pin
}

// NewRX480E4 ...
func NewRX480E4(out1, out2, out3, out4 uint8) *RX480E4 {
	r := &RX480E4{
		out1: rpio.Pin(out1),
		out2: rpio.Pin(out2),
		out3: rpio.Pin(out3),
		out4: rpio.Pin(out4),
	}
	r.out1.Input()
	r.out2.Input()
	r.out3.Input()
	r.out4.Input()
	return r
}

// PressA ...
func (r *RX480E4) PressA() bool {
	return r.out1.Read() == rpio.High
}

// PressB ...
func (r *RX480E4) PressB() bool {
	return r.out2.Read() == rpio.High
}

// PressC ...
func (r *RX480E4) PressC() bool {
	return r.out3.Read() == rpio.High
}

// PressD ...
func (r *RX480E4) PressD() bool {
	return r.out4.Read() == rpio.High
}
