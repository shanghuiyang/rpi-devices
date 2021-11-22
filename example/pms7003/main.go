package main

import (
	"fmt"
	"log"
	"time"

	"github.com/shanghuiyang/rpi-devices/dev"
)

const (
	devName = "/dev/ttyAMA0"
	baud    = 9600
)

func main() {
	p, err := dev.NewPMS7003(devName, baud)

	if err != nil {
		log.Printf("failed to new pms7003, error: %v", err)
		return
	}
	defer p.Close()

	for {
		time.Sleep(5000 * time.Millisecond)
		pm25, pm10, err := p.Get()
		if err != nil {
			log.Printf("failed to get data, error: %v", err)
			continue
		}
		fmt.Printf("pm2.5: %v, pm10:: %v\n", pm25, pm10)
	}
}
