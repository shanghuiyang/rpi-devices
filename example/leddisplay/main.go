package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/jakefau/rpi-devices/dev"
	"github.com/stianeikeland/go-rpio"
)

const (
	dioPin  = 9
	rclkPin = 10
	sclkPin = 11
)

func main() {

	if err := rpio.Open(); err != nil {
		log.Fatalf("failed to open rpio, error: %v", err)
		return
	}
	defer rpio.Close()

	d := dev.NewLedDisplay(dioPin, rclkPin, sclkPin)
	d.Open()
	for {
		fmt.Printf(">>input: ")
		reader := bufio.NewReader(os.Stdin)
		input, err := reader.ReadString('\n')
		if err != nil {
			log.Printf("invalid input, error: %v", err)
			break
		}
		input = strings.Trim(input, "\n")
		if input == "q!" {
			log.Printf("quit")
			break
		}
		d.Display(input)

	}
	d.Close()
}
