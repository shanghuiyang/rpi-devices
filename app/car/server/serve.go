package server

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/shanghuiyang/rpi-devices/util"
)

const logTag = "server"

func Start(cfg *Config) {
	s, err := newService(cfg)
	if err != nil {
		log.Panicf("failed to new service, error: %v", err)
	}

	if err := s.start(); err != nil {
		log.Panicf("[%v]failed to start server, error: %v", logTag, err)
	}

	log.Printf("[%v]service started", logTag)
	util.WaitQuit(func() {
		s.shutdown()
	})

	r := mux.NewRouter()
	routeAPIs(r, s)

	http.Handle("/", r)
	if err := http.ListenAndServe(cfg.Host, nil); err != nil {
		log.Panicf("[%v]failed to start http server, error: %v", logTag, err)
	}
	log.Printf("[%v]http server stop", logTag)
}

func routeAPIs(r *mux.Router, s *service) {
	// home
	r.HandleFunc("/", s.loadHomeHandler).Methods("GET")

	// car operation
	r.HandleFunc("/car/{op:[a-z]+}", s.opHandler).Methods("POST")

	// car turn an angle
	r.HandleFunc("/car/turn/{angle}", s.turnHandler).Methods("POST")

	// self-driving
	r.HandleFunc("/selfdriving/on", s.selfDrivingOnHandler).Methods("POST")
	r.HandleFunc("/selfdriving/off", s.selfDrivingOffHandler).Methods("POST")

	// self-tracking
	r.HandleFunc("/selftracking/on", s.selfTrackingOnHandler).Methods("POST")
	r.HandleFunc("/selftracking/off", s.selfTrackingOffHandler).Methods("POST")

	// speech-driving
	r.HandleFunc("/speechdriving/on", s.speechDrivingOnHandler).Methods("POST")
	r.HandleFunc("/speechdriving/off", s.speechDrivingOffHandler).Methods("POST")

	// self-nav
	r.HandleFunc("/selfnav/{lat}/{lon}", s.selfNavOnHandler).Methods("POST")
	r.HandleFunc("/selfnav/off", s.selfNavOffHandler).Methods("POST")

	// set volume
	r.HandleFunc("/volume/{v:[0-9]+}", s.setVolumeHandler).Methods("POST")

	// music
	r.HandleFunc("/music/on", s.musicOnHandler).Methods("POST")
	r.HandleFunc("/music/off", s.musicOffHandler).Methods("POST")
}
