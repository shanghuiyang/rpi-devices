/*
Encoder is a sensor used to count number.

Connect to Raspberry Pi:
 - vcc: any 3.3v or v5 pin
 - gnd: any gnd pin
 - out: any data pin

*/
package dev

import (
	"github.com/stianeikeland/go-rpio/v4"
)

// EncoderImp implements Encoder interface
type EncoderImp struct {
	pin rpio.Pin
}

// NewEncoderImp ...
func NewEncoderImp(pin uint8) *EncoderImp {
	e := &EncoderImp{
		pin: rpio.Pin(pin),
	}
	e.pin.Input()
	e.pin.PullDown()
	e.pin.Detect(rpio.NoEdge)
	return e
}

// Detected ...
func (e *EncoderImp) Detected() bool {
	return e.pin.EdgeDetected()
}

// Count1 ...
func (e *EncoderImp) Count1() int {
	if e.pin.EdgeDetected() {
		return 1
	}
	return 0
}

// Start ...
func (e *EncoderImp) Start() {
	e.pin.Detect(rpio.RiseEdge)
}

// Stop ...
func (e *EncoderImp) Stop() {
	e.pin.Detect(rpio.NoEdge)
}
