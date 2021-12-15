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
)

// const (
// 	devName = "/dev/ttyAMA0"
// 	baud    = 9600
// )

type (
	option func(s *sserver)
)

type sserver struct {
	thermohygrometer dev.Thermohygrometer
	pms7003          *dev.PMS7003
}

type temphumiResponse struct {
	Temp     float64 `json:"temp"`
	Humi     float64 `json:"humi"`
	ErrorMsg string  `json:"error_msg"`
}

type pm25Response struct {
	PM25     uint16 `json:"pm25"`
	ErrorMsg string `json:"error_msg"`
}

func main() {
	hdc, err := dev.NewHDC1080()
	if err != nil {
		log.Printf("[sensors]failed to new HDC1080")
		return
	}

	// p, err := dev.NewPMS7003(devName, baud)
	// if err != nil {
	// 	log.Printf("[sensors]failed to new PMS7003")
	// 	return
	// }

	s := newServer(
		withThermohygrometer(hdc),
		// withPMS7003(p),
	)
	if s == nil {
		log.Fatal("[sensors]failed to new sserver")
		return
	}

	if err := s.start(); err != nil {
		log.Printf("[sensors]failed to start sserver, error: %v", err)
		os.Exit(1)
	}
	os.Exit(0)
}

func withThermohygrometer(t dev.Thermohygrometer) option {
	return func(s *sserver) {
		s.thermohygrometer = t
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
	http.HandleFunc("/temphumi", s.temphumiHandler)
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

func (s *sserver) temphumiHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("[sensors]%v %v", r.Method, r.URL.Path)
	if s.thermohygrometer == nil {
		resp := &temphumiResponse{
			ErrorMsg: "invaild thermohygrometer sensor",
		}
		s.response(w, resp, http.StatusInternalServerError)
		return
	}

	t, h, err := s.thermohygrometer.TempHumidity()
	if err != nil {
		resp := &temphumiResponse{
			ErrorMsg: fmt.Sprintf("failed to get temp and humi, error: %v", err),
		}
		s.response(w, resp, http.StatusInternalServerError)
		return
	}

	resp := &temphumiResponse{
		Temp: t,
		Humi: h,
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
