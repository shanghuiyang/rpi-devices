package main

import (
	"fmt"
	"log"
	"time"

	"github.com/shanghuiyang/rpi-devices/dev"
)

var celsiusChar = []byte{0xDF}
var celsiusStr string = string(celsiusChar[:])

func main() {
	display, err := dev.NewLcdDisplay(16, 2)
	if err != nil {
		log.Printf("failed to create an oled, error: %v", err)
		return
	}
	defer display.Close()

	if err := display.On(); err != nil {
		log.Printf("failed to turn backlight on, error: %v", err)
		return
	}

	//  display following context on a 16x2 LCD.
	//
	//   0 _____________ 15
	// 0 |     27.3'C     |
	// 1 |    15:04:05    |
	//   +----------------+
	//
	text := fmt.Sprintf("27.3%sC", celsiusStr) // 27.3'C
	if err := display.Text(text, 5, 0); err != nil {
		log.Printf("failed to display time, error: %v", err)
		return
	}

	for i := 0; i < 10; i++ {
		t := time.Now().Format("15:04:05")
		if err := display.Text(t, 4, 1); err != nil {
			log.Printf("failed to display time, error: %v", err)
			break
		}
		time.Sleep(1 * time.Second)
	}
}
