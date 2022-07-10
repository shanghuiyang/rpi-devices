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
		if !rfp.Detected() {
			time.Sleep(100 * time.Millisecond)
			continue
		}
		log.Printf("detect pressure")
		time.Sleep(1 * time.Second)
	}
}
