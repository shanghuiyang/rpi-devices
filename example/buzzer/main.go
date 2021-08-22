package main

import (
	"fmt"
	"log"

	"github.com/shanghuiyang/rpi-devices/dev"
)

const (
	pin = 10
)

func main() {

	buz := dev.NewBuzzerImp(pin, true)

	var op string
	for {
		fmt.Printf(">>op: ")
		if n, err := fmt.Scanf("%s", &op); n != 1 || err != nil {
			log.Printf("invalid operator, error: %v", err)
			continue
		}
		switch op {
		case "on":
			buz.On()
		case "off":
			buz.Off()
		case "beep":
			buz.Beep(3, 300)
		case "q":
			log.Printf("quit\n")
			return
		default:
			fmt.Printf("invalid operator, should be: on, off, blink or q\n")
		}
	}
}
