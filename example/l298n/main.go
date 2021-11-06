package main

import (
	"fmt"
	"log"
	"time"

	"github.com/shanghuiyang/rpi-devices/dev"
)

const (
	in1 = 17
	in2 = 23
	in3 = 24
	in4 = 22
	ena = 18
	enb = 13
)

func main() {

	l298n := dev.NewL298N(in1, in2, in3, in4, ena, enb)
	l298n.SetSpeed(30)
	var op string
	for {
		fmt.Printf(">>op: ")
		if n, err := fmt.Scanf("%s", &op); n != 1 || err != nil {
			log.Printf("invalid operator, error: %v", err)
			continue
		}
		switch op {
		case "f":
			l298n.Forward()
		case "b":
			l298n.Backward()
		case "l":
			l298n.Left()
		case "r":
			l298n.Right()
		case "q":
			log.Printf("quit\n")
			return
		default:
			fmt.Printf("invalid operator, should be: f(forward), b(backward), l(left), r(right) or q(quit)\n")
		}
		time.Sleep(1 * time.Second)
		l298n.Stop()
	}
}
