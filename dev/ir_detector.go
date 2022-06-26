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
	out rpio.Pin
}

// NewIRDetector ...
func NewIRDetector(out uint8) *IRDetector {
	ir := &IRDetector{
		out: rpio.Pin(out),
	}
	ir.out.Input()
	return ir
}

// Detected ...
func (ir *IRDetector) Detected() bool {
	return ir.out.Read() == rpio.Low
}
