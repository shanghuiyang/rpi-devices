package main

import (
	"log"
	"time"

	"github.com/shanghuiyang/rpi-devices/dev"
)

const (
	do = 25
)

func main() {
	mq7 := dev.NewMQ7(do)
	for {
		if !mq7.Detected() {
			time.Sleep(100 * time.Millisecond)
			continue
		}
		log.Printf("detect CO")
		time.Sleep(10 * time.Second)
	}
}
