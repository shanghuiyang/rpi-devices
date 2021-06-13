/*
Package dev ...
SG90 is the driver of servo motor

Connect to Pi:
 - the red line:	any 5v pin
 - the brown line: 	any gnd pin
 - the yellow line:	any pwn pin(must be one of gpio 12, 13, 18, 19)
*/
package dev

import (
	"time"

	"github.com/jakefau/rpi-devices/util"
	"github.com/stianeikeland/go-rpio"
)

// SG90 ...
type SG90 struct {
	pin rpio.Pin
	rpi util.RpiModel
}

// NewSG90 ...
func NewSG90(pin uint8) *SG90 {
	s := &SG90{
		pin: rpio.Pin(pin),
		rpi: util.GetRpiModel(),
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
//         eye
//
func (s *SG90) Roll(angle int) {
	if angle < -90 || angle > 90 {
		return
	}
	duty := uint32(10.0 - float32(angle)/15.0)
	if s.rpi == util.Rpi4 {
		// Rpi4 uses a BCM 2711, which is different from the early rpi like rpi3, rpi2, rpiA and rpi0
		duty = uint32(26.5 - 39*float32(angle)/180)
	}
	s.pin.DutyCycle(uint32(duty), 100)
	time.Sleep(100 * time.Millisecond)
	s.pin.DutyCycle(0, 100)
}
