package main

import (
	"log"
	"time"

	"github.com/shanghuiyang/rpi-devices/dev"
)

const (
	d0 = 23
	d1 = 24
	d2 = 25
	d3 = 8

	ledPin = 21

	channelButonA = 3
	channelButonB = 2
	channelButonC = 1
	channelButonD = 0
)

func main() {
	r := dev.NewRX480E4(d0, d1, d2, d3)
	led := dev.NewLedImp(ledPin)
	
	defer func() {
		led.Off()
	}()

	ledOn := false
	chA := make(chan bool)
	chB := make(chan bool)
	chC := make(chan bool)
	chD := make(chan bool)

	go func(ch chan bool) {
		for {
			if r.Received(channelButonA) {
				ch <- true
				time.Sleep(1 * time.Second)
				continue
			}
			time.Sleep(20 * time.Millisecond)
		}
	}(chA)

	go func(ch chan bool) {
		for {
			if r.Received(channelButonB) {
				ch <- true
				time.Sleep(1 * time.Second)
				continue
			}
			time.Sleep(20 * time.Millisecond)
		}
	}(chB)

	go func(ch chan bool) {
		for {
			if r.Received(channelButonC) {
				ch <- true
				time.Sleep(1 * time.Second)
				continue
			}
			time.Sleep(20 * time.Millisecond)
		}
	}(chC)

	go func(ch chan bool) {
		for {
			if r.Received(channelButonD) {
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
