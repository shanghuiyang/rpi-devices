/*
sserver is a sensor server which provide data from all kinds of sensors
*/

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/shanghuiyang/rpi-devices/dev"
	"github.com/shanghuiyang/rpi-devices/util"
	"github.com/stianeikeland/go-rpio"
)

type (
	option func(s *sserver)
)

type sserver struct {
	ds18b20 *dev.DS18B20
	pms7003 *dev.PMS7003
}

type tempResponse struct {
	Temp     float32 `json:"temp"`
	ErrorMsg string  `json:"error_msg"`
}

type pm25Response struct {
	PM25     uint16 `json:"pm25"`
	ErrorMsg string `json:"error_msg"`
}

func main() {
	if err := rpio.Open(); err != nil {
		log.Fatalf("[sensors]failed to open rpio, error: %v", err)
		os.Exit(1)
	}
	defer rpio.Close()

	d := dev.NewDS18B20()
	if d == nil {
		log.Printf("[sensors]failed to new DS18B20")
		return
	}

	p := dev.NewPMS7003()
	if p == nil {
		log.Printf("[sensors]failed to new PMS7003")
		return
	}

	s := newServer(
		withDS18B20(d),
		withPMS7003(p),
	)
	if s == nil {
		log.Fatal("[sensors]failed to new sserver")
		return
	}

	util.WaitQuit(func() {
		rpio.Close()
	})
	if err := s.start(); err != nil {
		log.Printf("[sensors]failed to start sserver, error: %v", err)
		os.Exit(1)
	}
	os.Exit(0)
}

func withDS18B20(d *dev.DS18B20) option {
	return func(s *sserver) {
		s.ds18b20 = d
	}
}

func withPMS7003(p *dev.PMS7003) option {
	return func(s *sserver) {
		s.pms7003 = p
	}
}

func newServer(opts ...option) *sserver {
	s := &sserver{}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

func (s *sserver) start() error {
	log.Printf("[sensors]start service")
	http.HandleFunc("/temp", s.tempHandler)
	http.HandleFunc("/pm25", s.pm25Handler)
	if err := http.ListenAndServe(":8000", nil); err != nil {
		return err
	}
	return nil
}

func (s *sserver) response(w http.ResponseWriter, resp interface{}, statusCode int) error {
	w.WriteHeader(statusCode)
	data, err := json.Marshal(resp)
	if err != nil {
		log.Printf("[sensors]failed to marshal the response, error: %v", err)
		return err
	}
	if _, err := w.Write(data); err != nil {
		log.Printf("[sensors]failed to write data to http.ResponseWriter, error: %v", err)
		return err
	}
	return nil
}

func (s *sserver) tempHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("[sensors]%v %v", r.Method, r.URL.Path)
	if s.ds18b20 == nil {
		resp := &tempResponse{
			ErrorMsg: "invaild ds18b20 sensor",
		}
		s.response(w, resp, http.StatusInternalServerError)
		return
	}

	t, err := s.ds18b20.GetTemperature()
	if err != nil {
		resp := &tempResponse{
			ErrorMsg: fmt.Sprintf("failed to get temp, error: %v", err),
		}
		s.response(w, resp, http.StatusInternalServerError)
		return
	}

	resp := &tempResponse{
		Temp: t,
	}
	s.response(w, resp, http.StatusOK)
}

func (s *sserver) pm25Handler(w http.ResponseWriter, r *http.Request) {
	log.Printf("[sensors]%v %v", r.Method, r.URL.Path)
	if s.pms7003 == nil {
		resp := &pm25Response{
			ErrorMsg: "invaild pms7003 sensor",
		}
		s.response(w, resp, http.StatusInternalServerError)
		return
	}

	pm25, _, err := s.pms7003.Get()
	if err != nil {
		resp := &pm25Response{
			ErrorMsg: fmt.Sprintf("failed to get pm2.5, error: %v", err),
		}
		s.response(w, resp, http.StatusInternalServerError)
		return
	}

	resp := &pm25Response{
		PM25: pm25,
	}
	s.response(w, resp, http.StatusOK)
}
