/*
GY25 is an accelerometer used to detect yaw, pitch and roll angles.

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
  - VCC: any 3.3v or 5v pin
  - GND: any gnd pin
  - RXT: must connect to GPIO-14/TXD
  - TXD: must connect to GPIO-15/RXD
*/
package dev

import (
	"errors"
	// "fmt"
	"time"

	"github.com/tarm/serial"
)

const (
	gy25zBufSize   = 11
	gy25zDataHead  = 0x5A
	gy25zAngleData = 0x10
	gy25zDataLen   = 0x06
)

// GY25ZMode ...
type GY25ZMode []byte

var (
	// GY25ZQueryMode ...
	GY25ZQueryMode = GY25ZMode{0xA5, 0x56, 0x01, 0xFC}
	// GY25ZAutoMode ...
	GY25ZAutoMode = GY25ZMode{0xA5, 0x56, 0x02, 0xFD}
	// GY25ZAutoTextMode ...
	GY25ZAutoTextMode = GY25ZMode{0xA5, 0x53}
	// GY25ZCorrectionPitchAndRollMode ...
	GY25ZCorrectionPitchAndRollMode = GY25ZMode{0xA5, 0x54}
	// GY25ZCorrectionYawMode ...
	GY25ZCorrectionYawMode = GY25ZMode{0xA5 + 0x55}
)

// GY25Z implements Accelerometer interface
type GY25Z struct {
	port *serial.Port
	buf  [gy25zBufSize]byte
}

// NewGY25Z ...
func NewGY25Z(dev string, baud int) (*GY25Z, error) {
	cfg := &serial.Config{
		Name:        dev,
		Baud:        baud,
		ReadTimeout: 3 * time.Second,
	}
	port, err := serial.OpenPort(cfg)
	if err != nil {
		return nil, err
	}
	return &GY25Z{port: port}, nil
}

// SetMode ...
func (gy *GY25Z) SetMode(mode GY25ZMode) error {
	if err := gy.port.Flush(); err != nil {
		return err
	}

	n, err := gy.port.Write(mode)
	if n != 2 || err != nil {
		return err
	}
	return nil
}

// Angles ...
func (gy *GY25Z) Angles() (yaw, pitch, roll float64, err error) {
	if err := gy.port.Flush(); err != nil {
		return 0, 0, 0, err
	}

	a := 0
	for a < gy25zBufSize {
		n, err := gy.port.Read(gy.buf[a:])
		if err != nil {
			return 0, 0, 0, err
		}
		a += n
	}

	var data []byte
	for i := range gy.buf {
		if gy.buf[i] == gy25zDataHead && gy.buf[i+1] == gy25zDataHead && gy.buf[i+2] == gy25zAngleData && gy.buf[i+3] == gy25zDataLen {
			if i+4+gy25zDataLen < gy25zBufSize {
				data = gy.buf[i+4 : i+4+gy25zDataLen]
				break
			}
		}
	}

	if len(data) < gy25zDataLen {
		return 0, 0, 0, errors.New("unexpected data length")
	}

	r := (int16(data[0]) << 8) | int16(data[1])
	p := (int16(data[2]) << 8) | int16(data[3])
	y := (int16(data[4]) << 8) | int16(data[5])
	return float64(y) / 100, float64(p) / 100, float64(r) / 100, nil
}

// Close ...
func (gy *GY25Z) Close() error {
	return gy.port.Close()
}
