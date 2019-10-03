package devices

import (
	"github.com/stianeikeland/go-rpio"
)

// SG90 ...
type SG90 struct {
	pin rpio.Pin
}

// NewSG90 ...
func NewSG90(pin uint8) *SG90 {
	s := &SG90{
		pin: rpio.Pin(pin),
	}
	s.pin.Pwm()
	s.pin.Freq(50)
	s.pin.DutyCycle(0, 100)
	return s
}

// Roll ...
func (s *SG90) Roll(angle int) {
	duty := uint32(16.0 - float32(angle)/15.0)
	s.pin.DutyCycle(duty, 100)
}
