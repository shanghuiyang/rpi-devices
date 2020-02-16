package main

import (
	"log"
	"time"

	"github.com/shanghuiyang/rpi-devices/base"
	"github.com/shanghuiyang/rpi-devices/dev"
	"github.com/stianeikeland/go-rpio"
)

const (
	out1 = 1
	out2 = 2
	out3 = 4
	out4 = 5
)

func main() {
	if err := rpio.Open(); err != nil {
		log.Fatalf("failed to open rpio, error: %v", err)
		return
	}
	defer rpio.Close()

	r := dev.NewRX480E4(out1, out2, out3, out4)
	base.WaitQuit(func() {
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
