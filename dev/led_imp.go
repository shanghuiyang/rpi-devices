/*
LedImp is led light diodes module.

Connect to Raspberry Pi:
  - vcc(the longer pin) :  any data pin(~3.3v)
  - gnd(the shorter pin): any gnd pin
*/
package dev

import (
	"time"

	"github.com/stianeikeland/go-rpio/v4"
)

// LedImp implements Led interface
type LedImp struct {
	pin rpio.Pin
}

// NewLedImp ...
func NewLedImp(pin uint8) *LedImp {
	led := &LedImp{
		pin: rpio.Pin(pin),
	}
	led.pin.Output()
	led.pin.Low()
	return led
}

// On ...
func (led *LedImp) On() {
	led.pin.High()
}

// Off ...
func (led *LedImp) Off() {
	led.pin.Low()
}

// Blink is let led blink n time, interval Millisecond each time
func (led *LedImp) Blink(n int, intervalMs int) {
	d := time.Duration(intervalMs)
	for i := 0; i < n; i++ {
		led.On()
		delayMs(d)
		led.Off()
		delayMs(d)
	}
}

// Fade ...
func (led *LedImp) Fade(n uint8) {
	led.pin.Pwm()
	led.pin.Freq(64000)
	led.pin.DutyCycle(0, 32)
	for i := uint8(0); i < n; i++ {
		for j := uint32(0); j < 32; j++ { // increasing brightness
			led.pin.DutyCycle(j, 32)
			time.Sleep(time.Second / 32)
		}
		for j := uint32(32); j > 0; j-- { // decreasing brightness
			led.pin.DutyCycle(j, 32)
			time.Sleep(time.Second / 32)
		}
	}
	led.pin.Output()
	led.pin.Low()
}
