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
		log.Fatalf("failed to open rpio, error: %v", err)
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
			log.Printf("pressed A")
			light.turn()
		case <-chB:
			log.Printf("pressed B")
			light.turn()
		case <-chC:
			log.Printf("pressed C")
			light.turn()
		case <-chD:
			log.Printf("pressed D")
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
