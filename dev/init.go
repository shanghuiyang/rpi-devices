package dev

import (
	"log"

	"github.com/stianeikeland/go-rpio"
)

func init() {
	log.Print("init rpio in dev package")
	if err := rpio.Open(); err != nil {
		panic("failed to init rpio")
	}
}
