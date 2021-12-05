package main

import (
	"log"
	"time"

	"github.com/shanghuiyang/rpi-devices/dev"
)

const (
	devName = "/dev/ttyAMA0"
	baud    = 9600
)

var (
	tvVolumeDown     = []byte{0xA1, 0xF1, 0xB3, 0x4C, 0x81}
	tvChannelForward = []byte{0xA1, 0xF1, 0xB3, 0x4C, 0xCA}
)

func main() {
	ir, err := dev.NewIRCoder(devName, baud)
	if err != nil {
		log.Fatalf("new ircoder error: %v", err)
	}
	defer ir.Close()

	for i := 0; i < 5; i++ {
		time.Sleep(5 * time.Second)
		if err := ir.Send(tvVolumeDown); err != nil {
			log.Printf("error: %v", err)
			continue
		}
		log.Print("sent successfully")
	}
}
