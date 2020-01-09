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

	"github.com/shanghuiyang/rpi-devices/base"
	"github.com/tarm/serial"
)

const (
	logTagZE08CH2O = "ZE08CH2O"
	maxDeltaCH2O   = 0.06
)

var (
	mockCH2Os       = []float64{0.05, 0.08, 0.1, 0.07}
	mockCH2OArryIdx = -1
)

// ZE08CH2O ...
type ZE08CH2O struct {
	port     *serial.Port
	buf      [32]byte
	history  *base.History
	maxRetry int
}

// NewZE08CH2O ...
func NewZE08CH2O() *ZE08CH2O {
	p := &ZE08CH2O{
		history:  base.NewHistory(10),
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
			log.Printf("[%v]incorrect data len: %v, expected: 9", logTagZE08CH2O, a)
			continue
		}

		checksum := ^(p.buf[1] + p.buf[2] + p.buf[3] + p.buf[4] + p.buf[5] + p.buf[6] + p.buf[7]) + 1
		if checksum != p.buf[8] {
			log.Printf("[%v]checksum failure", logTagZE08CH2O)
			continue
		}
		ppm := (uint16(p.buf[4]) << 8) | uint16(p.buf[5])
		ch2o := float64(ppm) * 0.001228 // convert ppm to mg/m3

		if !p.checkDelta(ch2o) {
			log.Printf("[%v]check delta failed, discard current data. CH2O: %v mg/m3", logTagZE08CH2O, ch2o)
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

func (p *ZE08CH2O) checkDelta(ch2o float64) bool {
	avg, err := p.history.Avg()
	if err != nil {
		if err == base.ErrEmpty {
			p.history.Add(ch2o)
			return true
		}
		return false
	}

	passed := math.Abs(avg-ch2o) < maxDeltaCH2O
	if passed {
		p.history.Add(ch2o)
	}
	return passed
}

// Mock ...
func (p *ZE08CH2O) Mock() (float64, error) {
	mockCH2OArryIdx++
	if mockCH2OArryIdx == len(mockCH2Os) {
		mockCH2OArryIdx = 0
	}
	return mockCH2Os[mockCH2OArryIdx], nil
}
