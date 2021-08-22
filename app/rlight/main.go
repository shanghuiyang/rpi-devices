package main

import (
	"log"
	"time"

	"github.com/shanghuiyang/rpi-devices/dev"
	"github.com/shanghuiyang/rpi-devices/util"
)

const (
	d0 = 16
	d1 = 20
	d2 = 21
	d3 = 6

	ledPin = 26

	butonAchannel = 3
	butonBchannel = 2
	butonCchannel = 1
	butonDchannel = 0
)

var light *rlight

type rlight struct {
	led   dev.Led
	state bool // on of off
}

func main() {
	led := dev.NewLedImp(ledPin)
	light = &rlight{
		led:   led,
		state: false,
	}
	r := dev.NewRX480E4(d0, d1, d2, d3)

	util.WaitQuit(func() {
		led.Off()
	})

	for {
		if r.Received(butonAchannel) {
			log.Printf("[rlight]pressed A")
			go light.toggle()
			continue
		}
		if r.Received(butonBchannel) {
			log.Printf("[rlight]pressed B")
			go light.toggle()
			continue
		}
		if r.Received(butonCchannel) {
			log.Printf("[rlight]pressed C")
			go light.toggle()
			continue
		}
		if r.Received(butonDchannel) {
			log.Printf("[rlight]pressed D")
			go light.toggle()
			continue
		}
		time.Sleep(300 * time.Millisecond)
	}
}

func (r *rlight) toggle() {
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
