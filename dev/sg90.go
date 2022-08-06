/*
SG90 is servo motor which can roll angels from 0~180 degree.

Connect to Raspberry Pi:
 - the red line:	any 5v pin
 - the brown line: 	any gnd pin
 - the yellow line:	must be one of gpio 12, 13, 18 or 19 (pwn pins)
*/
package dev

import (
	"github.com/stianeikeland/go-rpio/v4"
)

// SG90 implements Motor interface
type SG90 struct {
	pin rpio.Pin
	rpi rpiModel
}

// NewSG90 ...
func NewSG90(pin uint8) *SG90 {
	sg := &SG90{
		pin: rpio.Pin(pin),
		rpi: getRpiModel(),
	}
	sg.pin.Pwm()
	sg.pin.Freq(50)
	sg.pin.DutyCycle(0, 100)
	return sg
}

// Roll ...
// angle: [-90, 90]
// angle < 0: roll anticlockwise
// angel = 0: ahead
// angle > 0: roll clockwise
// e.g.
//
//       -30  0   30
//         \  |  /
//          \ | /
//           \|/
//   -90 ---- * ---- 90
//         +-----+
//         |     |
//         |     | sg90
//         | (*) |
//         +-----+
//
func (sg *SG90) Roll(angle float64) {
	if angle < -90 || angle > 90 {
		return
	}
	duty := uint32(10.0 - angle/15.0)
	if sg.rpi == rpi4 {
		// Rpi4 uses a BCM 2711,
		// which is different from
		// the early rpi like rpi3, rpi2, rpiA and rpi0
		duty = uint32(23.5 - 0.155*float64(angle))
	}
	sg.pin.DutyCycle(uint32(duty), 100)
}

func (sg *SG90) SetSpeed(speed int) {
	// do notiong just implement Motor interface.
}
