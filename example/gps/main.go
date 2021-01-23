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
	g := dev.NewGPS(devName, baud)
	defer g.Close()

	for {
		time.Sleep(1 * time.Second)
		pt, err := g.Loc()
		if err != nil {
			log.Printf("failed, error: %v", err)
			continue
		}
		log.Printf("%v", pt)
	}
}
