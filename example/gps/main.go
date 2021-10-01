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
	gps, err := dev.NewNeo6mGPS(devName, baud)
	if err != nil {
		log.Printf("failed to create gps, error: %v", err)
		return
	}
	defer gps.Close()

	for {
		time.Sleep(1 * time.Second)
		pt, err := gps.Loc()
		if err != nil {
			log.Printf("failed, error: %v", err)
			continue
		}
		log.Printf("%v", pt)
	}
}
