/*
PumpImp is a driver for pump motor module.

Connect to Raspberry Pi:
  - positive(the longer pin): 	any data pin
  - negative(she shorter pin): 	any gnd pin
*/
package dev

import (
	"time"

	"github.com/stianeikeland/go-rpio/v4"
)

// PumpImp implements Pump interface
type PumpImp struct {
	pin rpio.Pin
}

// NewLedImp ...
func NewPumpImp(pin uint8) *PumpImp {
	p := &PumpImp{
		pin: rpio.Pin(pin),
	}
	p.pin.Output()
	p.pin.Low()
	return p
}

// On ...
func (p *PumpImp) On() {
	p.pin.High()
}

// Off ...
func (p *PumpImp) Off() {
	p.pin.Low()
}

// Blink is let led blink n time, interval Millisecond each time
func (p *PumpImp) Run(sec int) {
	p.On()
	delaySec(time.Duration(sec))
	p.Off()
}
