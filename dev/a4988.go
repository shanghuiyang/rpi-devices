/*
A4988 is a driver for A4988 driver board used to drive stepper motor such as NAME.

Connect to Pi:
 - vmot: 5v pin of rpi, or 5~35v power supply for motor
 - gnd : any gnd pin of rpi, or the gnd of power supply
 - vcc : any 5v pin
 - gnd : any gnd pin
 - step: any data pin
 - dir : any data pin
 - ms1 : any data pin
 - ms2 : any data pin
 - ms3 : any data pin

 NOTE:
 - This driver didn't implement reset, sleep and enable/disable.
 - You should typing [sleep] to [reset] pin to get the module always doesn't sleep.
 - And get [enable] tying nothing to get module always be enabled.
*/
package dev

import (
	"errors"

	"github.com/stianeikeland/go-rpio/v4"
)

// the degree per step for nema stepper
var degreePerStepForNema = map[StepperMode]float64{
	FullMode:      1.8,
	HalfMode:      0.9,
	QuarterMode:   0.45,
	EighthMode:    0.225,
	SixteenthMode: 0.1125,
}

// A4988 ...
type A4988 struct {
	step          rpio.Pin
	dir           rpio.Pin
	ms1, ms2, ms3 rpio.Pin
	mode          StepperMode
}

// NewA4988 ...
func NewA4988(step, dir, ms1, ms2, ms3 uint8) *A4988 {
	a := &A4988{
		step: rpio.Pin(step),
		dir:  rpio.Pin(dir),
		ms1:  rpio.Pin(ms1),
		ms2:  rpio.Pin(ms2),
		ms3:  rpio.Pin(ms3),
		mode: HalfMode,
	}
	a.step.Output()
	a.dir.Output()
	a.ms1.Output()
	a.ms2.Output()
	a.ms3.Output()
	a.step.Low()
	a.dir.Low()
	_ = a.SetMode(HalfMode)
	return a
}

// Step gets the motor rolls n steps.
// roll in clockwise direction if n > 0,
// or roll in counter-clockwise direction if n < 0,
// or motionless if n = 0.
func (a *A4988) Step(n int) {
	if n == 0 {
		return
	}
	if n > 0 {
		a.dir.High()
	} else {
		a.dir.Low()
		n = 0 - n
	}
	for i := 0; i < n; i++ {
		a.step.High()
		delayUs(500)
		a.step.Low()
		delayUs(500)
	}
}

// Roll gets the motor rolls angle degree.
// roll in clockwise direction if angle > 0,
// or roll in counter-clockwise direction if angle < 0,
// or motionless if angle = 0.
func (a *A4988) Roll(angle float64) {
	degree, ok := degreePerStepForNema[a.mode]
	if !ok {
		return
	}
	n := int(angle / degree)
	a.Step(n)
}

// SetMode sets the stepping mode
func (a *A4988) SetMode(mode StepperMode) error {
	switch mode {
	case FullMode:
		a.ms1.Low()
		a.ms2.Low()
		a.ms3.Low()
	case HalfMode:
		a.ms1.High()
		a.ms2.Low()
		a.ms3.Low()
	case QuarterMode:
		a.ms1.Low()
		a.ms2.High()
		a.ms3.Low()
	case EighthMode:
		a.ms1.High()
		a.ms2.High()
		a.ms3.Low()
	case SixteenthMode:
		a.ms1.High()
		a.ms2.High()
		a.ms3.High()
	default:
		return errors.New("invalid mode")
	}
	a.mode = mode
	return nil
}
