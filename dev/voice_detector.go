/*
VoiceDetector is an sensor used to detect voice.

Connect to Raspberry Pi:
 - vcc: any 3.3v or 5v pin
 - gnd: any gnd pin
 - out: any data pin

*/

package dev

import (
	"github.com/stianeikeland/go-rpio/v4"
)

// VoiceDetector implements Detector interface
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
