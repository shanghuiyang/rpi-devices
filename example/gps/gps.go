package main

import (
	"log"

	s "github.com/shanghuiyang/pi/devices"
)

func main() {
	g := s.NewGPS()
	pt, err := g.Loc()
	if err != nil {
		log.Printf("failed, error: %v", err)
		return
	}
	log.Printf("%v", pt)
	g.Close()
}
