/*
BuzzerImp is a buzzer module used to generate "beep, beep, ..." sound.

Connect to Pi:
 - vcc: any v3.3 pin
 - gnd: and gnd pin
 - i/o: any data pin
*/
package dev

import (
	"time"

	"github.com/stianeikeland/go-rpio"
)

// BuzzerImp implements Buzzer interface
type BuzzerImp struct {
	pin                rpio.Pin
	triggeredByHighTTL bool
}

// NewBuzzerImp ...
func NewBuzzerImp(pin uint8, triggeredByHighTTL bool) *BuzzerImp {
	b := &BuzzerImp{
		pin:                rpio.Pin(pin),
		triggeredByHighTTL: triggeredByHighTTL,
	}
	b.pin.Output()
	return b
}

// On ...
func (b *BuzzerImp) On() {
	if b.triggeredByHighTTL {
		b.pin.High()
		return
	}
	b.pin.Low()
}

// Off ...
func (b *BuzzerImp) Off() {
	if b.triggeredByHighTTL {
		b.pin.Low()
		return
	}
	b.pin.High()
}

// Beep beeps [n] times with an interval in [interval] millisecond
func (b *BuzzerImp) Beep(n int, intervalInMS int) {
	d := time.Duration(intervalInMS) * time.Millisecond
	for i := 0; i < n; i++ {
		b.On()
		time.Sleep(d)
		b.Off()
		time.Sleep(d)
	}
}
