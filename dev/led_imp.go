/*
LedImp is a led module.

Connect to Pi:
  - positive(the longer pin): 	any data pin
  - negative(she shorter pin): 	any gnd pin
*/
package dev

import (
	"time"

	"github.com/stianeikeland/go-rpio"
)

const (
	logTagLed = "led"
)

// LedImp implements Led interface
type LedImp struct {
	pin rpio.Pin
}

// NewLedImp ...
func NewLedImp(pin uint8) *LedImp {
	l := &LedImp{
		pin: rpio.Pin(pin),
	}
	l.pin.Output()
	return l
}

// On ...
func (l *LedImp) On() {
	l.pin.High()
}

// Off ...
func (l *LedImp) Off() {
	l.pin.Low()
}

// Blink is let led blink n time, interval Millisecond each time
func (l *LedImp) Blink(n int, interval int) {
	d := time.Duration(interval) * time.Millisecond
	for i := 0; i < n; i++ {
		l.On()
		time.Sleep(d)
		l.Off()
		time.Sleep(d)
	}
}

// Fade ...
func (l *LedImp) Fade(n uint8) {
	l.pin.Pwm()
	l.pin.Freq(64000)
	l.pin.DutyCycle(0, 32)
	for i := uint8(0); i < n; i++ {
		for j := uint32(0); j < 32; j++ { // increasing brightness
			l.pin.DutyCycle(j, 32)
			time.Sleep(time.Second / 32)
		}
		for j := uint32(32); j > 0; j-- { // decreasing brightness
			l.pin.DutyCycle(j, 32)
			time.Sleep(time.Second / 32)
		}
	}
	l.pin.Output()
	l.pin.Low()
}
