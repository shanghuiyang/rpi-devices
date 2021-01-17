package main

import (
	"log"
	"time"

	"github.com/shanghuiyang/rpi-devices/dev"
)

func main() {
	g := dev.NewGPS()
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
