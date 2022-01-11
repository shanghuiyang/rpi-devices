package main

import (
	"fmt"
	"log"

	"github.com/shanghuiyang/rpi-devices/dev"
)

const (
	ch1  = 0  // channel 1
	ch2  = 1  // channel 2
	ch3  = 2  // channel 2
	ch4  = 3  // channel 2
	pin1 = 26 // pin for channel 1
	pin2 = 26 // pin for channel 1
	pin3 = 26 // pin for channel 1
	pin4 = 26 // pin for channel 1
)

func main() {
	pins := []uint8{pin1, pin2, pin3, pin4}
	r := dev.NewRelayImp(pins)
	var op string
	for {
		fmt.Printf(">>op: ")
		if n, err := fmt.Scanf("%s", &op); n != 1 || err != nil {
			log.Printf("invalid operator, error: %v", err)
			continue
		}
		switch op {
		case "ch1on":
			r.On(ch1)
		case "ch2on":
			r.On(ch2)
		case "ch3on":
			r.On(ch3)
		case "ch4on":
			r.On(ch4)
		case "ch1off":
			r.Off(ch1)
		case "ch2off":
			r.Off(ch2)
		case "ch3off":
			r.Off(ch3)
		case "ch4off":
			r.Off(ch4)
		case "q":
			log.Printf("done\n")
			return
		default:
			fmt.Printf("invalid operator, should be: on, off or q\n")
		}
	}
}
