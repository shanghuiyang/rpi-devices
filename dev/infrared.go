package dev

import (
	"github.com/stianeikeland/go-rpio"
)

// InfraredDetector ...
type InfraredDetector struct {
	pin rpio.Pin
}

// NewInfraredDetector ...
func NewInfraredDetector(pin uint8) *InfraredDetector {
	i := &InfraredDetector{
		pin: rpio.Pin(pin),
	}
	i.pin.Input()
	return i
}

// Detected ...
func (i *InfraredDetector) Detected() bool {
	return i.pin.Read() == rpio.Low
}
