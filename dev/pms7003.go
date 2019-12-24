/*
PMS7003 is air quality sensor.

connecto to pi:
VCC: any 5v pin on pi
GND: any GND pin on pi
TXD: RXD
RXD: TXD
*/
package dev

import (
	"fmt"
	"log"
	"math"

	"github.com/tarm/serial"
)

const (
	logTagPMS7003 = "PMS7003"
)

var (
	mockPMs     = []uint16{50, 110, 150, 110, 50}
	mockArryIdx = -1
)

// PMS7003 ...
type PMS7003 struct {
	port     *serial.Port
	buf      [128]byte
	history  [10]uint16
	idx      uint8
	maxRetry int
}

// NewPMS7003 ...
func NewPMS7003() *PMS7003 {
	p := &PMS7003{
		maxRetry: 10,
	}
	if err := p.open(); err != nil {
		return nil
	}
	return p
}

// Get returns pm2.5 and pm10
func (p *PMS7003) Get() (uint16, uint16, error) {
	for i := 0; i < p.maxRetry; i++ {
		if err := p.port.Flush(); err != nil {
			return 0, 0, err
		}
		a := 0
		for a < 32 {
			n, err := p.port.Read(p.buf[a:])
			if err != nil {
				// try to reopen serial
				p.port.Close()
				if err := p.open(); err != nil {
					log.Printf("[%v]failed open serial, error: %v", logTagPMS7003, err)
				}
				return 0, 0, fmt.Errorf("error on read from port, error: %v. try to open serial again", err)
			}
			a += n
		}

		if a != 32 {
			return 0, 0, fmt.Errorf("incorrect data len for pm2.5, len: %v", a)
		}
		if p.buf[0] != 0x42 && p.buf[1] != 0x4d && p.buf[2] != 0 && p.buf[3] != 28 {
			return 0, 0, fmt.Errorf("bad data for pm2.5")
		}
		pm25 := (uint16(p.buf[6]) << 8) | uint16(p.buf[7])
		pm10 := (uint16(p.buf[8]) << 8) | uint16(p.buf[9])
		if !p.check(pm25) {
			log.Printf("[%v]check failed, discard current data. pm2.5: %v", logTagPMS7003, pm25)
			continue
		}
		return pm25, pm10, nil
	}
	return 0, 0, fmt.Errorf("failed to get pm2.5 and pm10")
}

// Close ...
func (p *PMS7003) Close() {
	p.port.Close()
}

func (p *PMS7003) open() error {
	c := &serial.Config{Name: "/dev/ttyAMA0", Baud: 9600}
	port, err := serial.OpenPort(c)
	if err != nil {
		return err
	}
	p.port = port
	return nil
}

func (p *PMS7003) check(pm25 uint16) bool {
	var (
		n   int
		sum float64
	)
	for _, pm25 := range p.history {
		if pm25 > 0 {
			sum += float64(pm25)
			n++
		}
	}
	if n == 0 {
		p.history[0] = pm25
		p.idx = 1
		return true
	}
	avg := sum / float64(n)
	passed := math.Abs(avg-float64(pm25)) < 200
	if passed {
		p.history[p.idx] = pm25
		p.idx++
		if p.idx > 9 {
			p.idx = 0
		}
	}
	return passed
}

// Mock ...
func (p *PMS7003) Mock() (uint16, uint16, error) {
	mockArryIdx++
	if mockArryIdx == len(mockPMs) {
		mockArryIdx = 0
	}
	return mockPMs[mockArryIdx], mockPMs[mockArryIdx], nil
}
