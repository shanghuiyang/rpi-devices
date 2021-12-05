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
	"math"
	"time"

	"github.com/tarm/serial"
)

const (
	maxDeltaCH2O        = 0.06
	ze08ch2oHistorySize = 10
)

// ZE08CH2O implements CH2OMeter interface
type ZE08CH2O struct {
	port     *serial.Port
	buf      [32]byte
	maxRetry int

	history []float64
	next    int
}

// NewZE08CH2O ...
func NewZE08CH2O() (*ZE08CH2O, error) {
	history := make([]float64, ze08ch2oHistorySize)
	for i := range history {
		history[i] = -1
	}

	ze := &ZE08CH2O{
		maxRetry: 10,
		history:  history,
		next:     0,
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
				if err := ze.port.Close(); err != nil {
					return 0, fmt.Errorf("close port error: %w", err)
				}
				if err := ze.open(); err != nil {
					return 0, fmt.Errorf("open port error: %w", err)
				}
			}
			a += n
		}

		if a != 9 {
			continue
		}

		checksum := ^(ze.buf[1] + ze.buf[2] + ze.buf[3] + ze.buf[4] + ze.buf[5] + ze.buf[6] + ze.buf[7]) + 1
		if checksum != ze.buf[8] {
			continue
		}
		ppm := (uint16(ze.buf[4]) << 8) | uint16(ze.buf[5])
		ch2o := float64(ppm) * 0.001228 // convert ppm to mg/m3

		if !ze.validate(ch2o) {
			continue
		}
		return ch2o, nil
	}
	return 0, fmt.Errorf("failed to get ch2o, arrived max retry times")
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

func (ze *ZE08CH2O) validate(ch2o float64) bool {
	if len(ze.history) == 0 {
		ze.history[0] = ch2o
		ze.next++
		return true
	}
	var sum float64
	n := 0
	for _, v := range ze.history {
		if v < 0 {
			break
		}
		sum += float64(v)
		n++
	}
	avg := sum / float64(n)
	if math.Abs(avg-float64(ch2o)) < maxDeltaCH2O {
		ze.history[ze.next] = ch2o
		ze.next++
		if ze.next >= ze08ch2oHistorySize {
			ze.next = 0
		}
		return true
	}
	return false
}
