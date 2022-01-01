package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/shanghuiyang/rpi-devices/dev"
)

const (
	dioPin  = 9
	rclkPin = 10
	sclkPin = 11
)

func main() {
	display := dev.NewTM1637Display(dioPin, rclkPin, sclkPin)
	display.On()
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
		display.Text(input, 0, 0)

	}
	display.Close()
}
