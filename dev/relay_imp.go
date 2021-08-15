/*
relay is an electrically operated switch module.

Connect to Pi:
 - vcc: any 5v pin
 - gnd: any gnd pin
 - in:  pin 26(gpio 7) or any data pin
 - on:  the outside device
 - com: the bettery

*/
package dev

import (
	"github.com/stianeikeland/go-rpio"
)

// RelayImp ...
type RelayImp struct {
	pins []rpio.Pin
}

// NewRelayImp ...
func NewRelayImp(pins []uint8) *RelayImp {
	p := make([]rpio.Pin, len(pins))
	for i, pin := range pins {
		p[i] = rpio.Pin(pin)
		p[i].Output()
	}
	r := &RelayImp{
		pins: p,
	}
	return r
}

// On ...
func (r *RelayImp) On(ch int) {
	if ch < 0 || ch >= len(r.pins) {
		return
	}
	r.pins[ch].High()
}

// Off ...
func (r *RelayImp) Off(ch int) {
	if ch < 0 || ch >= len(r.pins) {
		return
	}
	r.pins[ch].Low()
}
