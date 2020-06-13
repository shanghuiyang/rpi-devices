package main

import (
	"log"
	"time"

	"github.com/shanghuiyang/rpi-devices/base"
	"github.com/shanghuiyang/rpi-devices/dev"
	"github.com/stianeikeland/go-rpio"
)

const (
	d0 = 16
	d1 = 20
	d2 = 21
	d3 = 6

	ledPin = 26

	// use this rpio as 3.3v pin
	// if all 3.3v pins were used
	pin33v = 5
)

var light *rlight

type rlight struct {
	led   *dev.Led
	state bool // on of off
}

func main() {
	if err := rpio.Open(); err != nil {
		log.Fatalf("[rlight]failed to open rpio, error: %v", err)
		return
	}
	defer rpio.Close()

	p33v := rpio.Pin(pin33v)
	p33v.Output()
	p33v.High()

	led := dev.NewLed(ledPin)
	light = &rlight{
		led:   led,
		state: false,
	}
	r := dev.NewRX480E4(d0, d1, d2, d3)

	base.WaitQuit(func() {
		led.Off()
		rpio.Close()
	})

	for {
		if r.PressA() == true {
			log.Printf("[rlight]pressed A")
			go light.turn()
			continue
		}
		if r.PressB() == true {
			log.Printf("[rlight]pressed B")
			go light.turn()
			continue
		}
		if r.PressC() == true {
			log.Printf("[rlight]pressed C")
			go light.turn()
			continue
		}
		if r.PressD() == true {
			log.Printf("[rlight]pressed D")
			go light.turn()
			continue
		}
		time.Sleep(300 * time.Millisecond)
	}
}

func (r *rlight) turn() {
	if r.state {
		r.led.Off()
		r.state = false
		log.Printf("[rlight]light off")
	} else {
		r.led.On()
		r.state = true
		log.Printf("[rlight]light on")
	}
}
