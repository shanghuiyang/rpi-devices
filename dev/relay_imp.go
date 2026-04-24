/*
relay is an electrically operated switch module.

		           ___________________________
	   VCC/DC+ ---|                           |--- NC
	              |                           |
	   GND/DC- ---|          RELAY            |--- COM ------[+]Load[-]---------------------+
	              |                           |                                             |
	       IN  ---|                           |--- NO ------------------ [+]Bettery[-]------+
		           ---------------------------

Connect to Raspberry Pi:
  - vcc(dc+): any 5v pin
  - gnd(dc-): any gnd pin
  - in:  any data pin
  - on:  the outside device+
  - com: the bettery+
*/
package dev

import (
	"github.com/stianeikeland/go-rpio/v4"
)

// RelayImp ...
type RelayImp struct {
	pin rpio.Pin
}

// NewRelayImp ...
func NewRelayImp(pin uint8) *RelayImp {
	r := &RelayImp{
		pin: rpio.Pin(pin),
	}
	r.pin.Output()
	r.pin.Low()
	return r
}

// On ...
func (r *RelayImp) On() {
	r.pin.High()
}

// Off ...
func (r *RelayImp) Off() {
	r.pin.Low()
}
