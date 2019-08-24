package main

import (
	"fmt"
	"log"

	"github.com/shanghuiyang/pi/devices"
)

const (
	p7 = 7 // relay
)

func main() {
	r := devices.NewRelay(p7)
	go r.Start()

	var op string
	for {
		fmt.Printf(">>op: ")
		if n, err := fmt.Scanf("%s", &op); n != 1 || err != nil {
			log.Printf("invalid operator, error: %v", err)
			continue
		}
		switch op {
		case "on":
			devices.ChRelayOp <- devices.On
		case "off":
			devices.ChRelayOp <- devices.Off
		case "q":
			log.Printf("done\n")
			return
		default:
			fmt.Printf("invalid operator, should be: on, off or q\n")
		}
	}
}
