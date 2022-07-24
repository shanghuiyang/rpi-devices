package main

import (
	"log"
	"time"

	"github.com/shanghuiyang/rpi-devices/dev"
)

const (
	out = 25
)

func main() {
	ld := dev.NewLD2410(out)
	for {
		if !ld.Detected() {
			time.Sleep(1 * time.Second)
			continue
		}
		log.Printf("detected")
		time.Sleep(15 * time.Second)
	}
}
