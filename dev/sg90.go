/*
SG90 is servo motor which can roll angels from 0~180 degree.

Connect to Pi:
 - the red line:	any 5v pin
 - the brown line: 	any gnd pin
 - the yellow line:	any pwn pin(must be one of gpio 12, 13, 18, 19)
*/
package dev

import (
	"github.com/shanghuiyang/rpi-devices/util"
	"github.com/stianeikeland/go-rpio/v4"
)

// SG90 implements Servo interface
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
// angle < 0: roll anticlockwise
// angel = 0: ahead
// angle > 0: roll clockwise
// e.g.
//
//     -30  0   30
//       \  |  /
//        \ | /
//         \|/
//          *
//         eye
//
func (s *SG90) Roll(angle float64) {
	if angle < -90 || angle > 90 {
		return
	}
	duty := uint32(10.0 - angle/15.0)
	if s.rpi == util.Rpi4 {
		// Rpi4 uses a BCM 2711, which is different from the early rpi like rpi3, rpi2, rpiA and rpi0
		duty = uint32(26.5 - 39*float32(angle)/180)
	}
	s.pin.DutyCycle(uint32(duty), 100)
}

func (s *SG90) SetSpeed(speed int) {
	// Todo
}
