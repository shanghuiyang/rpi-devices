package main

import (
	"fmt"
	"log"

	"github.com/shanghuiyang/rpi-devices/dev"
)

const (
	ch1    = 0 // channel 1
	ch2    = 1 // channel 2
	ch1pin = 7 // pin for channel 1
	ch2pin = 8 // pin for channel 1
)

func main() {
	pins := []uint8{ch1pin, ch2pin}
	r := dev.NewRelayImp(pins)
	var op string
	for {
		fmt.Printf(">>op: ")
		if n, err := fmt.Scanf("%s", &op); n != 1 || err != nil {
			log.Printf("invalid operator, error: %v", err)
			continue
		}
		switch op {
		case "ch1 on":
			r.On(ch1)
		case "ch2 on":
			r.On(ch2)
		case "ch1 off":
			r.Off(ch1)
		case "ch2 off":
			r.Off(ch2)
		case "q":
			log.Printf("done\n")
			return
		default:
			fmt.Printf("invalid operator, should be: on, off or q\n")
		}
	}
}
