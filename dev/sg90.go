package dev

import (
	"time"

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
// angle: [-90, 90]
// angle < 0: left
// angel = 0: ahead
// angle > 0: right
// e.g.
//
//     -30  0   30
//       \  |  /
//        \ | /
//         \|/
//          *
//        camera
//
func (s *SG90) Roll(angle int) {
	if angle < -90 || angle > 90 {
		return
	}
	duty := uint32(10.0 - float32(angle)/15.0)
	s.pin.DutyCycle(duty, 100)
	time.Sleep(300 * time.Millisecond)
	s.pin.DutyCycle(0, 100)
}
