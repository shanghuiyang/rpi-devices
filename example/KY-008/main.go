package main

import (
	"github.com/shanghuiyang/rpi-devices/dev"
	"github.com/stianeikeland/go-rpio"
	"log"
	"time"
)

const (
	p14 = 14 // laser
)

func main(){
	if err := rpio.Open(); err != nil {
		log.Fatalf("failed to open rpio, error: %v", err)
		return
	}
	defer rpio.Close()

	laser := dev.NewLaser(p14)


	laser.On()
	time.Sleep(5 * time.Second)
	laser.Off()
	time.Sleep(5 * time.Second)


}