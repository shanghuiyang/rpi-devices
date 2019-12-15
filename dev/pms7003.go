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

	"github.com/tarm/serial"
)

const (
	logTagPMS7003 = "PMS7003"
)

var (
	pmBuf = make([]byte, 128)
)

// PMS7003 ...
type PMS7003 struct {
	port *serial.Port
}

// NewPMS7003 ...
func NewPMS7003() *PMS7003 {
	p := &PMS7003{}
	if err := p.open(); err != nil {
		return nil
	}
	return p
}

// Get returns pm2.5 and pm10
func (p *PMS7003) Get() (pm25, pm10 uint16, err error) {
	pm25 = 0
	pm10 = 0
	if err = p.port.Flush(); err != nil {
		return
	}
	a := 0
	for a < 32 {
		n, er := p.port.Read(pmBuf[a:])
		if er != nil {
			// try to reopen serial
			p.port.Close()
			if er := p.open(); er != nil {
				log.Printf("[%v]failed open serial, error: %v", logTagPMS7003, er)
			}
			err = fmt.Errorf("error on read from port, error: %v. try to open serial again", er)
			return
		}
		a += n
	}

	if a != 32 {
		err = fmt.Errorf("incorrect data len for pm2.5, len: %v", a)
		return
	}
	if pmBuf[0] != 0x42 && pmBuf[1] != 0x4d && pmBuf[2] != 0 && pmBuf[3] != 28 {
		err = fmt.Errorf("bad data for pm2.5")
		return
	}
	pm25 = (uint16(pmBuf[6]) << 8) | uint16(pmBuf[7])
	pm10 = (uint16(pmBuf[8]) << 8) | uint16(pmBuf[9])
	return
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
