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

func main() {
	if err := rpio.Open(); err != nil {
		log.Fatalf("failed to open rpio, error: %v", err)
		return
	}
	defer rpio.Close()

	r := dev.NewRX480E4(d0, d1, d2, d3)
	led := dev.NewLed(ledPin)
	base.WaitQuit(func() {
		led.Off()
		rpio.Close()
	})

	ledOn := false
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
			if ledOn {
				led.Off()
				ledOn = false
			} else {
				led.On()
				ledOn = true
			}
		case <-chB:
			log.Printf("pressed B")
		case <-chC:
			log.Printf("pressed C")
		case <-chD:
			log.Printf("pressed D")
		default:
			time.Sleep(20 * time.Millisecond)
		}
	}
}
