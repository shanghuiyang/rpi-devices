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
		detected := mq7.Detected()
		if detected {
			log.Printf("detected CO")
		} else {
			log.Printf("didn't detect CO")
		}
		time.Sleep(1 * time.Second)
	}
}
