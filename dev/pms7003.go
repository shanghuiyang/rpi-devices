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
	"time"

	"github.com/tarm/serial"
)

// PMS7003 ...
type PMS7003 struct {
	port  *serial.Port
	buf   [128]byte
	retry int
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
		port:  port,
		retry: 10,
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
				return 0, 0, err
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
		return pm25, pm10, nil
	}
	return 0, 0, errors.New("psm7003 is busy, please try agian later")
}

// Close ...
func (pms *PMS7003) Close() error {
	return pms.port.Close()
}
