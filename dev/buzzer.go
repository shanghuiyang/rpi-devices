package dev

import (
	"time"

	"github.com/stianeikeland/go-rpio"
)

// Buzzer ...
type Buzzer struct {
	pin rpio.Pin
}

// NewBuzzer ...
func NewBuzzer(pin int8) *Buzzer {
	b := &Buzzer{
		pin: rpio.Pin(pin),
	}
	b.pin.Output()
	return b
}

// Sound ...
func (b *Buzzer) Sound() {
	b.pin.High()
	time.Sleep(50 * time.Millisecond)
	b.pin.Low()
}
