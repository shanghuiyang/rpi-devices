package main

import (
	"log"
	"time"

	"github.com/shanghuiyang/rpi-devices/dev"
)

const (
	out = 18
)

func main() {
	ir := dev.NewIRDetector(out)
	for {
		if !ir.Detected() {
			time.Sleep(100 * time.Millisecond)
			continue
		}
		log.Printf("detected")
		time.Sleep(1 * time.Second)
	}
}
