/*
Package dev ...

ZE08CH2O is the driver of ZE08CH2O, an air quality sensor which can be used to detect PM2.5 and PM10.

Config Your Pi:
1. $ sudo vim /boot/config.txt
	add following new line:
	~~~~~~~~~~~~~~~~~
	enable_uart=1
	~~~~~~~~~~~~~~~~~
2. $ sudo vim /boot/cmdline.txt
	remove following contexts:
	~~~~~~~~~~~~~~~~~~~~~~~~~~
	console=serial0,115200
	~~~~~~~~~~~~~~~~~~~~~~~~~~
3. $ sudo reboot now
4. $ sudo cat /dev/ttyAMA0
	should see somethings output

Connect to Pi:
 - VCC: any 5v pin
 - GND: any gnd pin
 - TXD: must connect to pin 10(gpio 14) (RXD)

*/
package dev

import (
	"fmt"
	"log"
	"math"

	"github.com/tarm/serial"
)

const (
	logTagZE08CH2O = "ZE08CH2O"
)

var (
	mockCH2Os       = []float64{0.01, 0.05, 0.08, 0.1, 0.15}
	mockCH2OArryIdx = -1
)

// ZE08CH2O ...
type ZE08CH2O struct {
	port     *serial.Port
	buf      [128]byte
	history  [10]float64
	idx      uint8
	maxRetry int
}

// NewZE08CH2O ...
func NewZE08CH2O() *ZE08CH2O {
	p := &ZE08CH2O{
		maxRetry: 10,
	}
	if err := p.open(); err != nil {
		return nil
	}
	return p
}

// Get returns ch2o in mg/m3
func (p *ZE08CH2O) Get() (float64, error) {
	for i := 0; i < p.maxRetry; i++ {
		if err := p.port.Flush(); err != nil {
			return 0, err
		}
		a := 0
		for a < 9 {
			n, err := p.port.Read(p.buf[a:])
			if err != nil {
				// try to reopen serial
				p.port.Close()
				if err := p.open(); err != nil {
					log.Printf("[%v]failed open serial, error: %v", logTagZE08CH2O, err)
				}
				return 0, fmt.Errorf("error on read from port, error: %v. try to open serial again", err)
			}
			a += n
		}

		if a != 9 {
			return 0, fmt.Errorf("incorrect data len for ch2o, len: %v", a)
		}

		ppm := (uint16(p.buf[4]) << 8) | uint16(p.buf[5])
		ch2o := float64(ppm) * 0.001228 // convert ppm to mg/m3

		if !p.check(ch2o) {
			log.Printf("[%v]check failed, discard current data. CH2O: %v mg/m3", logTagZE08CH2O, ch2o)
			continue
		}
		return ch2o, nil
	}
	return 0, fmt.Errorf("failed to get ch2o")
}

// Close ...
func (p *ZE08CH2O) Close() {
	p.port.Close()
}

func (p *ZE08CH2O) open() error {
	c := &serial.Config{Name: "/dev/ttyAMA0", Baud: 9600}
	port, err := serial.OpenPort(c)
	if err != nil {
		return err
	}
	p.port = port
	return nil
}

func (p *ZE08CH2O) check(ch2o float64) bool {
	var (
		n   int
		sum float64
	)
	for _, h := range p.history {
		if h > 0 {
			sum += h
			n++
		}
	}
	if n == 0 {
		p.history[0] = ch2o
		p.idx = 1
		return true
	}
	avg := sum / float64(n)
	passed := math.Abs(avg-float64(ch2o)) < 0.1
	if passed {
		p.history[p.idx] = ch2o
		p.idx++
		if p.idx > 9 {
			p.idx = 0
		}
	}
	return passed
}

// Mock ...
func (p *ZE08CH2O) Mock() (float64, error) {
	mockCH2OArryIdx++
	if mockCH2OArryIdx == len(mockCH2Os) {
		mockArryIdx = 0
	}
	return mockCH2Os[mockCH2OArryIdx], nil
}
