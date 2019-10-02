package devices

import (
	"time"

	"github.com/stianeikeland/go-rpio"
)

const (
	logTagLed = "led"
)

// Led ...
type Led struct {
	pin  rpio.Pin
	isOn bool
}

// NewLed ...
func NewLed(pin uint8) *Led {
	l := &Led{
		pin:  rpio.Pin(pin),
		isOn: false,
	}
	l.pin.Output()
	return l
}

// On ...
func (l *Led) On() {
	if !l.isOn {
		l.pin.High()
		l.isOn = true
	}
}

// Off ...
func (l *Led) Off() {
	if l.isOn {
		l.pin.Low()
		l.isOn = false
	}
}

// Blink ...
func (l *Led) Blink(n uint8) {
	for i := uint8(0); i < n; i++ {
		l.On()
		time.Sleep(50 * time.Millisecond)
		l.Off()
		time.Sleep(50 * time.Millisecond)
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
