package main

import (
	"log"
	"time"

	"github.com/shanghuiyang/rpi-devices/dev"
)

const (
	devName = "/dev/ttyAMA0"
	baud    = 9600
)

func main() {
	zp, err := dev.NewZP16(devName, baud)
	if err != nil {
		log.Printf("failed to create gps, error: %v", err)
		return
	}
	defer zp.Close()

	for {
		time.Sleep(3 * time.Second)
		co, err := zp.CO()
		if err != nil {
			log.Printf("failed, error: %v", err)
			continue
		}
		log.Printf("%v", co)
	}
}
