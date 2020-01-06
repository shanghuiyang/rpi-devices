package main

import (
	"log"

	"github.com/shanghuiyang/rpi-devices/dev"
)

func main() {
	ch2o := dev.NewZE08CH2O()
	c, err := ch2o.Get()
	if err != nil {
		log.Printf("failed, error: %v", err)
		return
	}
	log.Printf("CH2O: %.4f mg/m3\n", c)

	ch2o.Close()
}
