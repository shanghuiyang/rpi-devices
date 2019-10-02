package main

import (
	"fmt"
	"log"

	dev "github.com/shanghuiyang/rpi-devices/devices"
	"github.com/stianeikeland/go-rpio"
)

const (
	p12 = 12 // led
)

func main() {
	if err := rpio.Open(); err != nil {
		log.Fatalf("failed to open rpio, error: %v", err)
		return
	}
	defer rpio.Close()

	led := dev.NewLed(p12)

	var op string
	for {
		fmt.Printf(">>op: ")
		if n, err := fmt.Scanf("%s", &op); n != 1 || err != nil {
			log.Printf("invalid operator, error: %v", err)
			continue
		}
		switch op {
		case "on":
			led.On()
		case "off":
			led.Off()
		case "blink":
			led.Blink(5)
		case "fade":
			led.Fade(3)
		case "q":
			log.Printf("quit\n")
			return
		default:
			fmt.Printf("invalid operator, should be: on, off, blink or q\n")
		}
	}
}
