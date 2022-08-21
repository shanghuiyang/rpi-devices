/*
BuzzerImp is a buzzer module used to generate "beep, beep, ..." sound.

Connect to Raspberry Pi for a 3-pin buzzer mobule:
 - vcc: any 3.3v pin
 - gnd: any gnd pin
 - i/o: any data pin

Connect to Raspberry Pi for a 2-pin buzzer mobule:
  - vcc(the longer pin):  any data pin(~3.3v)
  - gnd(the shorter pin): any gnd pin
*/

package dev

import (
	"time"

	"github.com/stianeikeland/go-rpio/v4"
)

// BuzzerImp implements Buzzer interface
type BuzzerImp struct {
	pin    rpio.Pin
	trigBy LogicLevel
}

// NewBuzzerImp ...
func NewBuzzerImp(pin uint8, trigBy LogicLevel) *BuzzerImp {
	b := &BuzzerImp{
		pin:    rpio.Pin(pin),
		trigBy: trigBy,
	}
	b.pin.Output()
	b.Off()
	return b
}

// On ...
func (b *BuzzerImp) On() {
	if b.trigBy == High {
		b.pin.High()
		return
	}
	b.pin.Low()
}

// Off ...
func (b *BuzzerImp) Off() {
	if b.trigBy == High {
		b.pin.Low()
		return
	}
	b.pin.High()
}

// Beep beeps [n] times with an interval in [interval] millisecond
func (b *BuzzerImp) Beep(n int, intervalMs int) {
	d := time.Duration(intervalMs)
	for i := 0; i < n; i++ {
		b.On()
		delayMs(d)
		b.Off()
		delayMs(d)
	}
}
