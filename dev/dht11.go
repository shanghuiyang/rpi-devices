package dev

import (
	"errors"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
	"time"
)

const (
	tFile = "/sys/bus/iio/devices/iio:device0/in_temp_input"
	hFile = "/sys/bus/iio/devices/iio:device0/in_humidityrelative_input"
)

const maxRetry = 50

// DHT11 ...
type DHT11 struct {
}

// NewDHT11 ...
func NewDHT11() *DHT11 {
	return &DHT11{}
}

// TempHumidity ...
func (d *DHT11) TempHumidity() (float64, float64, error) {
	var (
		t, h                 float64
		gotTemp, gotHumidity bool
	)

	for i := 0; i < maxRetry; i++ {
		if !gotTemp {
			if data, err := ioutil.ReadFile(tFile); err == nil {
				if t, err = d.parseData(data); err == nil {
					gotTemp = true
					log.Printf("got temp: %v", t)
				}
			} else {
				log.Printf("read tmep data failed: %v", err)
			}
		}

		if !gotHumidity {
			if data, err := ioutil.ReadFile(hFile); err == nil {
				if h, err = d.parseData(data); err == nil {
					gotHumidity = true
					log.Printf("got humidity: %v", h)
				}
			} else {
				log.Printf("read humidity data failed: %v", err)
			}
		}
		if gotTemp && gotHumidity {
			return t, h, nil
		}
		time.Sleep(1 * time.Second)
	}

	return 0, 0, errors.New("bad data")
}

func (d *DHT11) parseData(data []byte) (float64, error) {
	s := strings.Trim(string(data), " \t\n")
	log.Printf("data: %v", s)
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, err
	}
	return v / 1000.0, nil
}
