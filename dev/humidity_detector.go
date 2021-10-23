/*
HumidityDetector detects the humidity/soisture of soil.

Connect to Raspberry Pi:
 - vcc: any 3.3v or v5 pin
 - gnd: any gnd pin
 - out: any data pin

*/

package dev

import (
	"github.com/stianeikeland/go-rpio/v4"
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
