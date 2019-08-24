package main

import (
	"fmt"
	"log"

	"github.com/shanghuiyang/pi/devices"
)

const (
	p26 = 26 // led
)

func main() {
	led := devices.NewLed(p26)
	go led.Start()

	var op string
	for {
		fmt.Printf(">>op: ")
		if n, err := fmt.Scanf("%s", &op); n != 1 || err != nil {
			log.Printf("invalid operator, error: %v", err)
			continue
		}
		switch op {
		case "on":
			devices.ChLedOp <- devices.On
		case "off":
			devices.ChLedOp <- devices.Off
		case "blink":
			devices.ChLedOp <- devices.Blink
		case "q":
			log.Printf("done\n")
			return
		default:
			fmt.Printf("invalid operator, should be: on, off, blink or q\n")
		}
	}
}
