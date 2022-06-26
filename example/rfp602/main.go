package main

import (
	"log"
	"time"

	"github.com/shanghuiyang/rpi-devices/dev"
)

const (
	do = 23
)

func main() {
	rfp := dev.NewRFP602(do)
	for {
		detected := rfp.Detected()
		if detected {
			log.Printf("detected pressure")
		} else {
			log.Printf("didn't detect pressure")
		}
		time.Sleep(100 * time.Millisecond)
	}
}
