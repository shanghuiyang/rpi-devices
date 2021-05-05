/*
Package dev ...

Connect to Pi:
  - positive(R): 			   any data pin
  - positive(G): 			   any data pin
  - positive(B): 			   any data pin
  - negative(she shorter pin): 	any gnd pin
*/

package dev

import "github.com/stianeikeland/go-rpio"

const (
	logTagRGB = "rgb"
)

// RGBLed ...
type RGBLED struct {
	redPin   rpio.Pin
	greenPin rpio.Pin
	bluePin  rpio.Pin
}

// NewRGBLed ...
func NewRGBLed(redPin uint8, greenPin uint8, bluePin uint8) *RGBLED {
	l := &RGBLED{
		redPin:   rpio.Pin(redPin),
		greenPin: rpio.Pin(greenPin),
		bluePin:  rpio.Pin(bluePin),
	}
	l.redPin.Output()
	l.greenPin.Output()
	l.bluePin.Output()
	return l
}

func (r *RGBLED) RedOn() {
	r.redPin.High()
}

func (r *RGBLED) GreenOn() {
	r.greenPin.High()
}

func (r *RGBLED) BlueOn() {
	r.bluePin.High()
}

func (r *RGBLED) RedOff() {
	r.redPin.Low()
}

func (r *RGBLED) GreenOff() {
	r.greenPin.Low()
}

func (r *RGBLED) BlueOff() {
	r.bluePin.Low()
}
