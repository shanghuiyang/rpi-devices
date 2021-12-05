/*
DHT11 is a sensor used to meature temperature and humidity.

Config Raspberry Pi:
1. sudo vim /boot/config.txt
2. add following line to the end of config.txt
	--------------------------
	dtoverlay=dht11,gpiopin=4
	--------------------------

Connect to Raspberry Pi:
	VCC:	any 3.3v pin
	GND:	any gnd pin
	SIGNAL: must connect to GPIO-4
*/
package dev

import (
	"errors"
	"io/ioutil"
	"math"
	"strconv"
	"strings"
)

const (
	tFile = "/sys/bus/iio/devices/iio:device0/in_temp_input"
	hFile = "/sys/bus/iio/devices/iio:device0/in_humidityrelative_input"
)

// DHT11 implements Thermohygrometer interface
type DHT11 struct {
	tempHistory [10]float64
	humiHistory [10]float64
	tempIdx     uint8
	humiIdx     uint8
	maxRetry    int
}

// NewDHT11 ...
func NewDHT11() *DHT11 {
	return &DHT11{
		maxRetry: 50,
	}
}

// TempHumidity ...
func (d *DHT11) TempHumidity() (temp, humi float64, err error) {
	chTemp := make(chan float64)
	chHumi := make(chan float64)

	go func(ch chan float64) {
		for i := 0; i < d.maxRetry; i++ {
			data, err := ioutil.ReadFile(tFile)
			if err != nil {
				continue
			}
			t, err := d.parseData(data)
			if err != nil {
				continue
			}
			if !d.checkTemp(t) {
				continue
			}
			ch <- t
			return
		}
		ch <- -999
	}(chTemp)

	go func(ch chan float64) {
		for i := 0; i < d.maxRetry; i++ {
			data, err := ioutil.ReadFile(hFile)
			if err != nil {
				continue
			}
			h, err := d.parseData(data)
			if err != nil {
				continue
			}
			if !d.checkHumi(h) {
				continue
			}
			ch <- h
			return
		}
		ch <- -999
	}(chHumi)

	t := <-chTemp
	h := <-chHumi
	if t == -999 || h == -999 {
		return t, h, errors.New("dht11 isn't ready")
	}
	return t, h, nil
}

func (d *DHT11) parseData(data []byte) (float64, error) {
	s := strings.Trim(string(data), " \t\n")
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, err
	}
	if v == 0 {
		return 0, errors.New("dht11 isn't ready")
	}
	return v / 1000.0, nil
}

func (d *DHT11) checkTemp(temp float64) bool {
	var (
		n   int
		sum float64
	)
	for _, t := range d.tempHistory {
		if t > 0 {
			sum += float64(t)
			n++
		}
	}
	if n == 0 {
		d.tempHistory[0] = temp
		d.tempIdx = 1
		return true
	}
	avg := sum / float64(n)
	passed := math.Abs(avg-temp) < 10
	if passed {
		d.tempHistory[d.tempIdx] = temp
		d.tempIdx++
		if d.tempIdx > 9 {
			d.tempIdx = 0
		}
	}
	return passed
}

func (d *DHT11) checkHumi(humi float64) bool {
	var (
		n   int
		sum float64
	)
	for _, h := range d.humiHistory {
		if h > 0 {
			sum += float64(h)
			n++
		}
	}
	if n == 0 {
		d.humiHistory[0] = humi
		d.humiIdx = 1
		return true
	}
	avg := sum / float64(n)
	passed := math.Abs(avg-humi) < 20
	if passed {
		d.humiHistory[d.humiIdx] = humi
		d.humiIdx++
		if d.humiIdx > 9 {
			d.humiIdx = 0
		}
	}
	return passed
}
