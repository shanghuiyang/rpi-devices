/*
IRDetector is a sensor used to detected infrared ray.

Connect to Pi:
 - vcc: any 3.3v pin
 - gnd: any gnd pin
 - out: any data pin

*/
package dev

import (
	"github.com/stianeikeland/go-rpio/v4"
)

// IRDetector implements Detector interface
type IRDetector struct {
	pin rpio.Pin
}

// NewIRDetector ...
func NewIRDetector(pin uint8) *IRDetector {
	ir := &IRDetector{
		pin: rpio.Pin(pin),
	}
	ir.pin.Input()
	return ir
}

// Detected ...
func (ir *IRDetector) Detected() bool {
	return ir.pin.Read() == rpio.Low
}
