package main

import (
	"log"
	"time"

	"github.com/shanghuiyang/rpi-devices/dev"
)

const (
	pin = 18
)

func main() {
	ir := dev.NewIRDetector(pin)
	for {
		detectedObj := ir.Detected()
		if detectedObj {
			log.Printf("detected an object")
		} else {
			log.Printf("didn't detect any objects")
		}
		time.Sleep(1 * time.Second)
	}
}
