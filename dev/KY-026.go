package dev

import (
	"github.com/stianeikeland/go-rpio"
)

type KY026 struct {
	pin rpio.Pin
}

func NewKY026(pin uint8) *KY026{
	f := &KY026{
		pin: rpio.Pin(pin),
	}
	f.pin.Input()
	return f
}

func (f KY026) Detected() bool {
	return f.pin.Read() == rpio.High
}