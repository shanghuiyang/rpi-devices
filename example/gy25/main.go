package main

import (
	"fmt"
	"log"
	"time"

	"github.com/shanghuiyang/rpi-devices/dev"
)

const (
	devName = "/dev/ttyUSB0"
	baud    = 115200
)

func main() {
	g, err := dev.NewGY25(devName, baud)
	if err != nil {
		log.Fatalf("new gy25 error: %v", err)
	}
	defer g.Close()

	if err := g.SetMode(dev.GY25AutoMode); err != nil {
		log.Printf("failed to set mode, error: %v", err)
		return
	}

	for {
		time.Sleep(100 * time.Millisecond)
		yaw, pitch, roll, err := g.Angles()
		if err != nil {
			log.Printf("failed to gent angles, error: %v", err)
			continue
		}
		fmt.Printf("yaw: %v, pitch: %v, roll: %v\n", yaw, pitch, roll)
	}
}
