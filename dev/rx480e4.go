/*
Package dev ...

Connect to Pi:
 - +v: any v3.3 or 5v pin
 - gnd: any gnd pin
 - d0: any data pin
 - d1: any data pin
 - d2: any data pin
 - d3: any data pin
 - vt: null

*/
package dev

import (
	"github.com/stianeikeland/go-rpio"
)

// RX480E4 ...
type RX480E4 struct {
	d0 rpio.Pin
	d1 rpio.Pin
	d2 rpio.Pin
	d3 rpio.Pin
}

// NewRX480E4 ...
func NewRX480E4(d0, d1, d2, d3 uint8) *RX480E4 {
	r := &RX480E4{
		d0: rpio.Pin(d0),
		d1: rpio.Pin(d1),
		d2: rpio.Pin(d2),
		d3: rpio.Pin(d3),
	}
	r.d0.Input()
	r.d1.Input()
	r.d2.Input()
	r.d3.Input()
	r.d0.PullDown()
	r.d1.PullDown()
	r.d2.PullDown()
	r.d3.PullDown()
	r.d0.Detect(rpio.RiseEdge)
	r.d1.Detect(rpio.RiseEdge)
	r.d2.Detect(rpio.RiseEdge)
	r.d3.Detect(rpio.RiseEdge)
	return r
}

// PressA ...
func (r *RX480E4) PressA() bool {
	return r.d3.EdgeDetected()
}

// PressB ...
func (r *RX480E4) PressB() bool {
	return r.d2.EdgeDetected()
}

// PressC ...
func (r *RX480E4) PressC() bool {
	return r.d1.EdgeDetected()
}

// PressD ...
func (r *RX480E4) PressD() bool {
	return r.d0.EdgeDetected()
}
