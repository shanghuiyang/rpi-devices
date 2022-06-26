/*
RFP-602 is a sensor used to detected pressure.

Connect to Pi:
 - vcc: any 3.3v pin
 - gnd: any gnd pin
 -  do: any data pin
 -  ao: (don't support currently)

*/
package dev

import (
	"github.com/stianeikeland/go-rpio/v4"
)

// RFP602 implements Detector interface
type RFP602 struct {
	do rpio.Pin
	ao rpio.Pin
}

// NewRFP602 ...
func NewRFP602(do uint8) *RFP602 {
	rfp := &RFP602{
		do: rpio.Pin(do),
	}
	rfp.do.Input()
	return rfp
}

// Detected ...
func (rfp *RFP602) Detected() bool {
	return rfp.do.Read() == rpio.Low
}
