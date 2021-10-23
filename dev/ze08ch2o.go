/*
ZE08CH2O is a sensor used to detect CH2O.

Config Raspberry Pi:
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

Connect to Raspberry Pi:
 - VCC: any 5v pin
 - GND: any gnd pin
 - TXD: must connect to GPIO-15 (RXD)

*/
package dev

import (
	"fmt"
	"log"
	"time"

	"math"

	"github.com/shanghuiyang/rpi-devices/util"
	"github.com/tarm/serial"
)

const (
	maxDeltaCH2O = 0.06
)

// ZE08CH2O implements CH2OMeter interface
type ZE08CH2O struct {
	port     *serial.Port
	buf      [32]byte
	history  *util.History
	maxRetry int
}

// NewZE08CH2O ...
func NewZE08CH2O() *ZE08CH2O {
	p := &ZE08CH2O{
		history:  util.NewHistory(10),
		maxRetry: 10,
	}
	if err := p.open(); err != nil {
		return nil
	}
	return p
}

// Get returns ch2o in mg/m3
func (p *ZE08CH2O) Value() (float64, error) {
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
					log.Printf("[ze08ch2o]failed open serial, error: %v", err)
				}
				return 0, fmt.Errorf("error on read from port, error: %v. try to open serial again", err)
			}
			a += n
		}

		if a != 9 {
			log.Printf("[ze08ch2o]incorrect data len: %v, expected: 9", a)
			continue
		}

		checksum := ^(p.buf[1] + p.buf[2] + p.buf[3] + p.buf[4] + p.buf[5] + p.buf[6] + p.buf[7]) + 1
		if checksum != p.buf[8] {
			log.Printf("[ze08ch2o]checksum failure")
			continue
		}
		ppm := (uint16(p.buf[4]) << 8) | uint16(p.buf[5])
		ch2o := float64(ppm) * 0.001228 // convert ppm to mg/m3

		if !p.checkDelta(ch2o) {
			log.Printf("[ze08ch2o]check delta failed, discard current data. CH2O: %v mg/m3", ch2o)
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
	c := &serial.Config{
		Name:        "/dev/ttyAMA0",
		Baud:        9600,
		ReadTimeout: 5 * time.Second,
	}
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
		if err == util.ErrEmpty {
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
