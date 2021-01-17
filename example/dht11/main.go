package main

import (
	"log"

	"github.com/shanghuiyang/rpi-devices/dev"
)

func main() {
	dht := dev.NewDHT11()
	t, h, err := dht.TempHumidity()
	if err != nil {
		log.Printf("failed, error: %v", err)
		return
	}
	log.Printf("t = %.0f, h = %.0f%%", t, h)
}
