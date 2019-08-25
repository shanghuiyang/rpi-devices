package main

import (
	"log"

	"github.com/shanghuiyang/rpi-devices/devices"
)

func main() {
	g := devices.NewGPS()
	pt, err := g.Loc()
	if err != nil {
		log.Printf("failed, error: %v", err)
		return
	}
	log.Printf("%v", pt)
	g.Close()
}
