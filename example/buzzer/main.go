package main

import (
	"github.com/shanghuiyang/rpi-devices/dev"
	"github.com/shanghuiyang/rpi-devices/util"
	"github.com/stianeikeland/go-rpio"
	"log"
	"time"
)

const buttonPin = 8
const buzzerPin = 26

func main(){

	if err := rpio.Open(); err != nil {
		log.Fatalf("failed to open rpio, error: %v", err)
		return
	}
	defer rpio.Close()

	buzz := dev.NewBuzzer(buzzerPin)
	btn := dev.NewButton(buttonPin)

	util.WaitQuit(func() {
		rpio.Close()
	})

	on := false
	for {
		pressed := btn.Pressed()
		if pressed {
			log.Printf("the button was pressed")
			if on {
				on = false
				buzz.Off()
			} else {
				buzz.On()
				on = true
			}
		}
		time.Sleep(300 * time.Millisecond)
	}

}
