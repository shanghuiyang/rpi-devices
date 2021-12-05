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

	"github.com/tarm/serial"
)

const (
	maxDeltaCH2O = 0.06
)

// ZE08CH2O implements CH2OMeter interface
type ZE08CH2O struct {
	port     *serial.Port
	buf      [32]byte
	history  *history
	maxRetry int
}

// NewZE08CH2O ...
func NewZE08CH2O() (*ZE08CH2O, error) {
	ze := &ZE08CH2O{
		history:  newHistory(10),
		maxRetry: 10,
	}
	if err := ze.open(); err != nil {
		return nil, err
	}
	return ze, nil
}

// Get returns ch2o in mg/m3
func (ze *ZE08CH2O) Value() (float64, error) {
	for i := 0; i < ze.maxRetry; i++ {
		if err := ze.port.Flush(); err != nil {
			return 0, err
		}
		a := 0
		for a < 9 {
			n, err := ze.port.Read(ze.buf[a:])
			if err != nil {
				// try to reopen serial
				ze.port.Close()
				if err := ze.open(); err != nil {
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

		checksum := ^(ze.buf[1] + ze.buf[2] + ze.buf[3] + ze.buf[4] + ze.buf[5] + ze.buf[6] + ze.buf[7]) + 1
		if checksum != ze.buf[8] {
			log.Printf("[ze08ch2o]checksum failure")
			continue
		}
		ppm := (uint16(ze.buf[4]) << 8) | uint16(ze.buf[5])
		ch2o := float64(ppm) * 0.001228 // convert ppm to mg/m3

		if !ze.checkDelta(ch2o) {
			log.Printf("[ze08ch2o]check delta failed, discard current data. CH2O: %v mg/m3", ch2o)
			continue
		}
		return ch2o, nil
	}
	return 0, fmt.Errorf("failed to get ch2o")
}

// Close ...
func (ze *ZE08CH2O) Close() error {
	return ze.port.Close()
}

func (ze *ZE08CH2O) open() error {
	c := &serial.Config{
		Name:        "/dev/ttyAMA0",
		Baud:        9600,
		ReadTimeout: 5 * time.Second,
	}
	port, err := serial.OpenPort(c)
	if err != nil {
		return err
	}
	ze.port = port
	return nil
}

func (ze *ZE08CH2O) checkDelta(ch2o float64) bool {
	avg, err := ze.history.Avg()
	if err != nil {
		if err == errEmpty {
			ze.history.Add(ch2o)
			return true
		}
		return false
	}

	passed := math.Abs(avg-ch2o) < maxDeltaCH2O
	if passed {
		ze.history.Add(ch2o)
	}
	return passed
}
