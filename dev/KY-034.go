package dev

import (
	"github.com/stianeikeland/go-rpio"
	"time"
)

type sevenColorLed struct {
	sig rpio.Pin
}

func NewSevenColorLed(pin uint) *sevenColorLed{
	led := &sevenColorLed{
		sig: rpio.Pin(pin),
	}
	led.sig.Output()
	return led
}

func (led *sevenColorLed) On(){
	led.sig.High()
}

func (led *sevenColorLed) Off(){
	led.sig.Low()
}

func (led *sevenColorLed) Blink(ttl time.Duration, times int){
	var i int
	for i < times {
		led.On()
		time.Sleep(ttl)
		led.Off()
		time.Sleep(ttl)
	}
}
