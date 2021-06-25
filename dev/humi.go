package dev

import (
	"github.com/stianeikeland/go-rpio"
)

// Humi ...
type Humi struct {
	pin rpio.Pin
}

// NewHumi ...
func NewHumi(pin uint8) *Humi {
	h := &Humi{
		pin: rpio.Pin(pin),
	}
	h.pin.Input()
	return h
}

// Detected ...
func (h *Humi) Detected() bool {
	return h.pin.Read() == rpio.Low
}
