package dev

import (
	"fmt"
	"github.com/stianeikeland/go-rpio"
)

type tap struct {
	pin rpio.Pin
}

func NewTapModule(pinNumber uint8)(*tap){
	m := &tap{
		pin: rpio.Pin(pinNumber),
	}
	m.pin.Output()
	return m
}

func (t tap) Tapped(threshold int) bool {
	if threshold > 1023 {
		threshold = 1023
	}
	fmt.Println(t.pin.Read())
	return t.pin.Read() == rpio.High
}

func (t tap) TappedStrength() rpio.State{
	return t.pin.Read()
}