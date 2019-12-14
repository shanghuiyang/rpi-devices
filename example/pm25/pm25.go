package main

import (
	"log"
	"time"

	"github.com/shanghuiyang/rpi-devices/dev"
)

func main() {
	p := dev.NewPMS7003()
	for {
		pm25, err := p.PM25()
		if err != nil {
			log.Printf("failed, error: %v", err)
			return
		}
		log.Printf("%v\n", pm25)
		time.Sleep(5 * time.Second)
	}
	p.Close()
}
