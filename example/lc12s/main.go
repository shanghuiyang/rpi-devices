package main

import (
	"fmt"
	"log"

	"github.com/shanghuiyang/rpi-devices/dev"
	"github.com/stianeikeland/go-rpio"
)

const (
	devName = "/dev/ttyAMA0"
	baud    = 9600
	csPin   = 2
)

func main() {
	fmt.Printf("1. I'm a sender.\n2. I'm a receiver.\n>>")

	var role string
	if n, err := fmt.Scanf("%s", &role); n != 1 || err != nil {
		log.Printf("invalid input, please input 1 or 2, error: %v", err)
		return
	}

	if err := rpio.Open(); err != nil {
		log.Fatalf("failed to open rpio, error: %v", err)
		return
	}
	defer rpio.Close()

	l, err := dev.NewLC12S(devName, baud, csPin) // receiver
	if err != nil {
		log.Fatalf("failed to new LC12S, error: %v", err)
		return
	}
	defer l.Close()
	l.Wakeup()
	switch role {
	case "1":
		sender(l)
	case "2":
		receiver(l)
	default:
		log.Printf("invalid input, please input 1 or 2")
	}
}

func sender(l *dev.LC12S) {
	for {
		fmt.Printf(">>send: ")

		var msg string
		if n, err := fmt.Scanf("%s", &msg); n != 1 || err != nil {
			log.Printf("invalid input, please input something, error: %v", err)
			continue
		}
		data := []byte(msg)
		if err := l.Send(data); err != nil {
			log.Printf("failed to send data, error: %v", err)
			continue
		}
		log.Printf("sent: %v", msg)
	}
}

func receiver(l *dev.LC12S) {
	for {
		data, err := l.Receive()
		if err != nil {
			log.Printf("failed to receive data, error: %v", err)
			continue
		}
		if len(data) == 0 {
			log.Printf("received noting")
			continue
		}
		log.Printf("received: %s, %v", data, data)
	}
}
