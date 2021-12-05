/*
BYJ2848 is a step-motor with 4-phase and 5-wire.

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
	"time"

	"github.com/stianeikeland/go-rpio/v4"
)

const (
	angleEachStep = 0.087
)

var (
	clockwise = [4][4]uint8{
		{1, 0, 0, 0},
		{0, 1, 0, 0},
		{0, 0, 1, 0},
		{0, 0, 0, 1},
	}

	anticlockwise = [4][4]uint8{
		{0, 0, 0, 1},
		{0, 0, 1, 0},
		{0, 1, 0, 0},
		{1, 0, 0, 0},
	}
)

// BYJ2848 implements Motor interface
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

func (byj *BYJ2848) Roll(angle float64) {
	var matrix [4][4]uint8
	if angle > 0 {
		matrix = clockwise
	} else {
		matrix = anticlockwise
		angle = angle * (-1)
	}
	n := int(angle / angleEachStep / 8.0)
	for i := 0; i < n; i++ {
		for j := 0; j < 4; j++ {
			for k := 0; k < 4; k++ {
				if matrix[j][k] == 1 {
					byj.pins[k].High()
				} else {
					byj.pins[k].Low()
				}
			}
			time.Sleep(2 * time.Millisecond)
		}
	}
}

// SetSpeed ...
func (byj *BYJ2848) SetSpeed() {
	// do noting just implement Motor interface.
}