package dev

import (
	"github.com/stianeikeland/go-rpio"
)

// VoiceDetector ...
type VoiceDetector struct {
	pin rpio.Pin
}

// NewVoiceDetector ...
func NewVoiceDetector(pin uint8) *VoiceDetector {
	v := &VoiceDetector{
		pin: rpio.Pin(pin),
	}
	v.pin.Input()
	return v
}

// Detected ...
func (v *VoiceDetector) Detected() bool {
	return v.pin.Read() == rpio.Low
}
