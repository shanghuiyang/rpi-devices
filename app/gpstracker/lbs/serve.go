package lbs

import (
	"log"
)

func Start(cfg *Config) {
	s, err := newService(cfg)
	if err != nil {
		log.Panicf("failed to new service, error: %v", err)
	}
	if err := s.start(); err != nil {
		log.Panicf("failed to start server, error: %v", err)
	}
}
