/*
Package dev ...

Connect to Pi:
 - vcc: any 3.3v or v5 pin
 - gnd: any gnd pin
 - out: any data pin

*/
package dev

import (
	"github.com/stianeikeland/go-rpio"
)

// Encoder ...
type Encoder struct {
	pin rpio.Pin
}

// NewEncoder ...
func NewEncoder(pin uint8) *Encoder {
	e := &Encoder{
		pin: rpio.Pin(pin),
	}
	e.pin.Input()
	e.pin.PullDown()
	e.pin.Detect(rpio.NoEdge)
	return e
}

// Detected ...
func (e *Encoder) Detected() bool {
	return e.pin.EdgeDetected()
}

// Count1 ...
func (e *Encoder) Count1() int {
	if e.pin.EdgeDetected() {
		return 1
	}
	return 0
}

// Start ...
func (e *Encoder) Start() {
	e.pin.Detect(rpio.RiseEdge)
}

// Stop ...
func (e *Encoder) Stop() {
	e.pin.Detect(rpio.NoEdge)
}
