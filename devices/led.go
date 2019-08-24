package devices

import (
	"log"
	"time"

	"github.com/stianeikeland/go-rpio"
)

const (
	logTagLed = "led"
)

var (
	// ChLedOp ...
	ChLedOp = make(chan Operator)
)

// Led ...
type Led struct {
	pin  rpio.Pin
	isOn bool
}

// NewLed ...
func NewLed(pin uint8) *Led {
	if err := rpio.Open(); err != nil {
		return nil
	}
	l := &Led{
		pin:  rpio.Pin(pin),
		isOn: false,
	}
	l.pin.Output()
	return l
}

// Start ...
func (l *Led) Start() {
	defer l.Close()

	log.Printf("[%v]start working", logTagLed)
	for {
		op := <-ChLedOp
		switch op {
		case Off:
			l.Off()
		case On:
			l.On()
		case Blink:
			l.Blink(5)
		default:
			// do nothing
		}
	}
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

// Close ...
func (l *Led) Close() {
	rpio.Close()
}
