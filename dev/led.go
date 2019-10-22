package dev

import (
	"time"

	"github.com/stianeikeland/go-rpio"
)

const (
	logTagLed = "led"
)

// Led ...
type Led struct {
	pin rpio.Pin
}

// NewLed ...
func NewLed(pin uint8) *Led {
	l := &Led{
		pin: rpio.Pin(pin),
	}
	l.pin.Output()
	return l
}

// On ...
func (l *Led) On() {
	l.pin.High()
}

// Off ...
func (l *Led) Off() {
	l.pin.Low()
}

// Blink ...
func (l *Led) Blink() {
	for {
		l.On()
		time.Sleep(1 * time.Second)
		l.Off()
		time.Sleep(1 * time.Second)
	}
}

// BlinkN is let led blink n time, interval Millisecond each time
func (l *Led) BlinkN(n int, interval int) {
	for i := 0; i < n; i++ {
		l.On()
		time.Sleep(time.Duration(interval) * time.Millisecond)
		l.Off()
		time.Sleep(time.Duration(interval) * time.Millisecond)
	}
}

// Fade ...
func (l *Led) Fade(n uint8) {
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
