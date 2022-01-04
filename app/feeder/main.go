package main

import (
	"fmt"
	"log"
	"time"

	"github.com/shanghuiyang/rpi-devices/dev"
	"github.com/shanghuiyang/rpi-devices/iot"
)

const (
	in1 = 12
	in2 = 16
	in3 = 20
	in4 = 21

	feedAt      = "10:00"
	onenetToken = "your_onenet_token"
	onenetAPI   = "http://api.heclouds.com/devices/540381180/datapoints"
)

func main() {
	cfg := &iot.Config{
		Token: onenetToken,
		API:   onenetAPI,
	}
	onenet := iot.NewOnenet(cfg)
	motor := dev.NewBYJ2848(in1, in2, in3, in4)

	var h, m int
	if n, err := fmt.Sscanf(feedAt, "%d:%d", &h, &m); n != 2 || err != nil {
		log.Fatalf("parse feed time error: %v", err)
	}
	log.Printf("feed at: %02d:%02d", h, m)

	total := 0
	for {
		now := time.Now()
		if now.Hour() == h && now.Minute() == m {
			motor.Roll(360)
			go push(onenet)
			total++
			log.Printf("fed, total: %v", total)
		}
		time.Sleep(time.Minute)
	}
}

func push(cloud iot.Cloud) {
	v := &iot.Value{
		Device: "feeder",
		Value:  1,
	}
	if err := cloud.Push(v); err != nil {
		log.Printf("push to clould error: %v", err)
		return
	}

	time.Sleep(10 * time.Second)
	v = &iot.Value{
		Device: "feeder",
		Value:  0,
	}
	if err := cloud.Push(v); err != nil {
		log.Printf("push to clould error: %v", err)
		return
	}
	log.Printf("push to cloud successfully")
}
