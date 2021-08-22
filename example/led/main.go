package main

import (
	"fmt"
	"log"

	"github.com/shanghuiyang/rpi-devices/dev"
)

const (
	p12 = 26 // led
)

func main() {

	led := dev.NewLedImp(p12)

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
			led.Blink(5, 100)
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
