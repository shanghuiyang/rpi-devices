/*
LD2410 is a sensor used to detected human existing or not.

Connect to Pi:
 - vcc: any 5v pin
 - gnd: any gnd pin
 - out: any data pin

*/
package dev

import (
	"github.com/stianeikeland/go-rpio/v4"
)

// LD2410 implements Detector interface
type LD2410 struct {
	out rpio.Pin
}

// NewLD2410 ...
func NewLD2410(out uint8) *LD2410 {
	ld := &LD2410{
		out: rpio.Pin(out),
	}
	ld.out.Input()
	return ld
}

// Detected ...
func (ld *LD2410) Detected() bool {
	return ld.out.Read() == rpio.High
}
