/*
PMS7003 is an air quality sensor which can be used to measure PM2.5 and PM10.

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
 - RX: must connect to GPIO-14/TXD
 - TX: must connect to GPIO-15/RXD

*/
package dev

import (
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/tarm/serial"
)

const (
	maxDeltaPM25 = 150
)

var (
	mockData = []uint16{50, 110, 150, 110, 50}
	mockIdx  = 0
)

// PMS7003 ...
type PMS7003 struct {
	port    *serial.Port
	history *history
	buf     [128]byte
	retry   int
}

// NewPMS7003 ...
func NewPMS7003(dev string, baud int) (*PMS7003, error) {
	cfg := &serial.Config{
		Name:        dev,
		Baud:        baud,
		ReadTimeout: 5 * time.Second,
	}
	port, err := serial.OpenPort(cfg)
	if err != nil {
		return nil, err
	}
	return &PMS7003{
		port:    port,
		history: newHistory(10),
		retry:   10,
	}, nil
}

// Get returns pm2.5 and pm10 in ug/m3
func (pms *PMS7003) Get() (uint16, uint16, error) {
	for i := 0; i < pms.retry; i++ {
		if err := pms.port.Flush(); err != nil {
			return 0, 0, err
		}
		a := 0
		for a < 32 {
			n, err := pms.port.Read(pms.buf[a:])
			if err != nil {
				return 0, 0, fmt.Errorf("error on read from port, error: %v. try to open serial again", err)
			}
			a += n
		}

		if a != 32 {
			continue
		}
		if pms.buf[0] != 0x42 && pms.buf[1] != 0x4d && pms.buf[2] != 0 && pms.buf[3] != 28 {
			continue
		}
		checksum := uint16(0)
		for i := 0; i < 29; i++ {
			checksum += uint16(pms.buf[i])
		}
		if checksum != (uint16(pms.buf[30])<<8)|uint16(pms.buf[31]) {
			continue
		}

		pm25 := (uint16(pms.buf[6]) << 8) | uint16(pms.buf[7])
		pm10 := (uint16(pms.buf[8]) << 8) | uint16(pms.buf[9])
		if !pms.checkDelta(pm25) {
			continue
		}
		return pm25, pm10, nil
	}
	return 0, 0, fmt.Errorf("psm7003 is invalid currently")
}

// MockGet mocks Get()
func (pms *PMS7003) MockGet() (uint16, uint16, error) {
	n := len(mockData)
	if n == 0 {
		return 0, 0, errors.New("without data")
	}
	if mockIdx >= n {
		mockIdx = 0
	}
	pm25 := mockData[mockIdx]
	pm10 := mockData[mockIdx]
	mockIdx++

	return pm25, pm10, nil
}

// Close ...
func (pms *PMS7003) Close() error {
	return pms.port.Close()
}

func (pms *PMS7003) checkDelta(pm25 uint16) bool {
	avg, err := pms.history.Avg()
	if err != nil {
		if err == errEmpty {
			pms.history.Add(pm25)
			return true
		}
		return false
	}

	passed := math.Abs(avg-float64(pm25)) < maxDeltaPM25
	if passed {
		pms.history.Add(pm25)
	}
	return passed
}
