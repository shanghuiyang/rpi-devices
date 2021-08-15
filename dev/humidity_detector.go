/*
HumidityDetector is a sensor used to deteched the humidity.
*/

package dev

import (
	"github.com/stianeikeland/go-rpio"
)

// HumidityDetector implements Detector interface
type HumidityDetector struct {
	pin rpio.Pin
}

// NewHumidityDetector ...
func NewHumidityDetector(pin uint8) *HumidityDetector {
	h := &HumidityDetector{
		pin: rpio.Pin(pin),
	}
	h.pin.Input()
	return h
}

// Detected ...
func (h *HumidityDetector) Detected() bool {
	return h.pin.Read() == rpio.Low
}
