package main

import (
	"log"
	"time"

	"github.com/shanghuiyang/rpi-devices/base"
	"github.com/shanghuiyang/rpi-devices/dev"
	"github.com/stianeikeland/go-rpio"
)

const (
	d0 = 23
	d1 = 24
	d2 = 25
	d3 = 8

	ledPin = 21
)

var light *rlight

type rlight struct {
	led   *dev.Led
	state bool // on of off
}

func main() {
	if err := rpio.Open(); err != nil {
		log.Fatalf("failed to open rpio, error: %v", err)
		return
	}
	defer rpio.Close()

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

	chA := make(chan bool)
	chB := make(chan bool)
	chC := make(chan bool)
	chD := make(chan bool)
	go func(ch chan bool) {
		for {
			if r.PressA() == true {
				ch <- true
				time.Sleep(1 * time.Second)
				continue
			}
			time.Sleep(20 * time.Millisecond)
		}
	}(chA)

	go func(ch chan bool) {
		for {
			if r.PressB() == true {
				ch <- true
				time.Sleep(1 * time.Second)
				continue
			}
			time.Sleep(20 * time.Millisecond)
		}
	}(chB)

	go func(ch chan bool) {
		for {
			if r.PressC() == true {
				ch <- true
				time.Sleep(1 * time.Second)
				continue
			}
			time.Sleep(20 * time.Millisecond)
		}
	}(chC)

	go func(ch chan bool) {
		for {
			if r.PressD() == true {
				ch <- true
				time.Sleep(1 * time.Second)
				continue
			}
			time.Sleep(20 * time.Millisecond)
		}
	}(chD)

	for {
		select {
		case <-chA:
			light.turn()
		case <-chB:
			light.turn()
		case <-chC:
			light.turn()
		case <-chD:
			light.turn()
		default:
			time.Sleep(20 * time.Millisecond)
		}
	}
}

func (r *rlight) turn() {
	if r.state {
		r.led.Off()
		r.state = false
		log.Printf("light off")
	} else {
		r.led.On()
		r.state = true
		log.Printf("light on")
	}
}
