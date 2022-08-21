/*
PumpImp is a driver for ~3.3v pump motor module.

Connect to Raspberry Pi:
  - vcc(red line)  : any data pin(~3.3v)
  - gnd(black line): any gnd pin
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

// Run lets the pump keep running in sec time
func (p *PumpImp) Run(sec int) {
	p.On()
	delaySec(time.Duration(sec))
	p.Off()
}
