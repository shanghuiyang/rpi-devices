package dev

import (
	"github.com/stianeikeland/go-rpio"
)

// Relay ...
type Relay struct {
	pin  rpio.Pin
	isOn bool
}

// NewRelay ...
func NewRelay(pin uint8) *Relay {
	r := &Relay{
		pin:  rpio.Pin(pin),
		isOn: false,
	}
	r.pin.Output()
	return r
}

// On ...
func (r *Relay) On() {
	if !r.isOn {
		r.pin.High()
		r.isOn = true
	}
}

// Off ...
func (r *Relay) Off() {
	if r.isOn {
		r.pin.Low()
		r.isOn = false
	}
}
