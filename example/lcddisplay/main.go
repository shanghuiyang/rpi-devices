package main

import (
	"log"
	"time"

	"github.com/shanghuiyang/rpi-devices/dev"
)

func main() {
	lcd, err := dev.NewLcdDisplay(16, 2)
	if err != nil {
		log.Printf("failed to create an oled, error: %v", err)
		return
	}
	defer lcd.Close()

	if err := lcd.BackLightOn(); err != nil {
		log.Printf("failed to turn backlight on, error: %v", err)
		return
	}
	if err := lcd.Display(3, 0, "hello world"); err != nil {
		log.Printf("failed to display time, error: %v", err)
		return
	}

	for i := 0; i < 10; i++ {
		t := time.Now().Format("15:04:05")
		if err := lcd.Display(4, 1, t); err != nil {
			log.Printf("failed to display time, error: %v", err)
			break
		}
		time.Sleep(1 * time.Second)
	}
}
