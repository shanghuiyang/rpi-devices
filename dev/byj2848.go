/*
BYJ2848 is a driver for stepper motor using ULN2003 driver board with 4-phase and 5-wire.

Connect to Pi:
 - vcc: any 5v pin
 - gnd: any gnd pin
 - in1:	any data pin
 - in2: any data pin
 - in3: any data pin
 - in4: any data pin

*/
package dev

import (
	"github.com/stianeikeland/go-rpio/v4"
)

const (
	degreePerStepForBYJ2848 = float64(0.703125) // 360/512
)

var (
	clockwise = [4][4]uint8{
		{1, 0, 0, 0},
		{0, 1, 0, 0},
		{0, 0, 1, 0},
		{0, 0, 0, 1},
	}

	cclockwise = [4][4]uint8{
		{0, 0, 0, 1},
		{0, 0, 1, 0},
		{0, 1, 0, 0},
		{1, 0, 0, 0},
	}
)

// BYJ2848 implements StepperMotor interface
type BYJ2848 struct {
	pins [4]rpio.Pin
}

// NewBYJ2848 ...
func NewBYJ2848(in1, in2, in3, in4 uint8) *BYJ2848 {
	byj := &BYJ2848{
		pins: [4]rpio.Pin{
			rpio.Pin(in1),
			rpio.Pin(in2),
			rpio.Pin(in3),
			rpio.Pin(in4),
		},
	}
	for i := 0; i < 4; i++ {
		byj.pins[i].Output()
		byj.pins[i].Low()
	}
	return byj
}

// Step gets the motor rolls n steps.
// roll in clockwise direction if n > 0,
// or roll in counter-clockwise direction if n < 0,
// or motionless if n = 0.
func (byj *BYJ2848) Step(n int) {
	matrix := clockwise
	if n < 0 {
		matrix = cclockwise
		n = 0 - n
	}

	for i := 0; i < n; i++ {
		for j := 0; j < 4; j++ {
			for k := 0; k < 4; k++ {
				if matrix[j][k] == 1 {
					byj.pins[k].High()
				} else {
					byj.pins[k].Low()
				}
			}
			delayMs(2)
		}
	}
	byj.reset()
}

// Roll gets the motor rolls angle degree.
// roll in clockwise direction if angle > 0,
// or roll in counter-clockwise direction if angle < 0,
// or motionless if angle = 0.
func (byj *BYJ2848) Roll(angle float64) {
	n := int(angle / degreePerStepForBYJ2848)
	byj.Step(n)
}

// SetMode sets the stepping mode.
// Please NOTE only FullMode is supported currently, and FullMode is used by default.
func (byj *BYJ2848) SetMode(mode StepperMode) error {
	return nil
}

func (byj *BYJ2848) reset() {
	for i := 0; i < 4; i++ {
		byj.pins[i].Low()
	}
}
