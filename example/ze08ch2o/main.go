package main

import (
	"log"

	"github.com/shanghuiyang/rpi-devices/dev"
)

func main() {
	z, err := dev.NewZE08CH2O()
	if err != nil {
		log.Fatalf("new ze08ch2o error: %v", err)
	}

	c, err := z.Value()
	if err != nil {
		log.Printf("get ch2o error: %v", err)
		return
	}
	log.Printf("CH2O: %.4f mg/m3\n", c)

	z.Close()
}
