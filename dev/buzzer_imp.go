/*
BuzzerImp is a buzzer module used to generate "beep, beep, ..." sound.

Connect to Raspberry Pi for a 3-pin buzzer mobule:
 - vcc: any 3.3v pin
 - gnd: any gnd pin
 - i/o: any data pin

Connect to Raspberry Pi for a 2-pin buzzer mobule:
 - port-1: any 3.3v pin
 - port-2: any data pin
*/

package dev

import (
	"time"

	"github.com/stianeikeland/go-rpio/v4"
)

// BuzzerImp implements Buzzer interface
type BuzzerImp struct {
	pin        rpio.Pin
	trigByHigh bool
}

// NewBuzzerImp ...
func NewBuzzerImp(pin uint8, trigByHigh bool) *BuzzerImp {
	b := &BuzzerImp{
		pin:        rpio.Pin(pin),
		trigByHigh: trigByHigh,
	}
	b.pin.Output()
	return b
}

// On ...
func (b *BuzzerImp) On() {
	if b.trigByHigh {
		b.pin.High()
		return
	}
	b.pin.Low()
}

// Off ...
func (b *BuzzerImp) Off() {
	if b.trigByHigh {
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
