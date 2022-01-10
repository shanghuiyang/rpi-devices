package main

import (
	"log"
	"time"

	"github.com/shanghuiyang/rpi-devices/dev"
)

const (
	pin = 6
)

func main() {
	e := dev.NewEncoderImp(pin)
	e.Start()
	defer e.Stop()

	count := 0
	for {
		count += e.Count1()
		log.Printf("%v", count)
		time.Sleep(100 * time.Millisecond)
	}
}
