package main

import (
	"log"

	"github.com/shanghuiyang/rpi-devices/dev"
	"github.com/stianeikeland/go-rpio"
)

func main() {
	if err := rpio.Open(); err != nil {
		log.Fatalf("failed to open rpio, error: %v", err)
		return
	}
	defer rpio.Close()

	l := dev.NewLC12S(26)
	defer l.Close()

	l.Wakeup()
	for {
		data, err := l.Read()
		if err != nil {
			log.Printf("failed to read data, error: %v", err)
			continue
		}
		if len(data) == 0 {
			log.Printf("received noting")
			continue
		}
		log.Printf("received: %v", data)
	}
}
