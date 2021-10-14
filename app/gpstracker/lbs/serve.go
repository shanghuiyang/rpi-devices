package lbs

import (
	"log"
)

const logTag = "lbs"

func Start(cfg *Config) {
	s, err := newService(cfg)
	if err != nil {
		log.Panicf("failed to new service, error: %v", err)
	}
	if err := s.start(); err != nil {
		log.Panicf("[%v]failed to start server, error: %v", logTag, err)
	}
}
