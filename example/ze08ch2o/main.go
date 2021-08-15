package main

import (
	"log"

	"github.com/shanghuiyang/rpi-devices/dev"
)

func main() {
	z := dev.NewZE08CH2O()
	c, err := z.Value()
	if err != nil {
		log.Printf("failed, error: %v", err)
		return
	}
	log.Printf("CH2O: %.4f mg/m3\n", c)

	z.Close()
}
