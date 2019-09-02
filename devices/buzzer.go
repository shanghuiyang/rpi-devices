package devices

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
	if err := rpio.Open(); err != nil {
		return nil
	}
	b := &Buzzer{
		pin: rpio.Pin(pin),
	}
	b.pin.Output()
	return b
}

// Whistle ...
func (b *Buzzer) Whistle() {
	b.pin.High()
	time.Sleep(100 * time.Millisecond)
	b.pin.Low()
}

// Close ...
func (b *Buzzer) Close() {
	rpio.Close()
}
