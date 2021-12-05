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
	"time"

	"github.com/tarm/serial"
)

const (
	datalen  = 8
	bufsize  = 16
	datahead = 0xAA
	datatail = 0x55
)

// GY25Mode ...
type GY25Mode []byte

var (
	// GY25QueryMode ...
	GY25QueryMode = GY25Mode{0xA5, 0x51}
	// GY25AutoMode ...
	GY25AutoMode = GY25Mode{0xA5, 0x52}
	// GY25AutoTextMode ...
	GY25AutoTextMode = GY25Mode{0xA5, 0x53}
	// GY25CorrectionPitchAndRollMode ...
	GY25CorrectionPitchAndRollMode = GY25Mode{0xA5, 0x54}
	// GY25CorrectionYawMode ...
	GY25CorrectionYawMode = GY25Mode{0xA5 + 0x55}
)

// GY25 implements Accelerometer interface
type GY25 struct {
	port *serial.Port
	buf  [bufsize]byte
}

// NewGY25 ...
func NewGY25(dev string, baud int) (*GY25, error) {
	cfg := &serial.Config{
		Name:        dev,
		Baud:        baud,
		ReadTimeout: 3 * time.Second,
	}
	port, err := serial.OpenPort(cfg)
	if err != nil {
		return nil, err
	}
	return &GY25{port: port}, nil
}

// SetMode ...
func (gy *GY25) SetMode(mode GY25Mode) error {
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
func (gy *GY25) Angles() (yaw, pitch, roll float64, err error) {
	if err := gy.port.Flush(); err != nil {
		return 0, 0, 0, err
	}

	a := 0
	for a < 16 {
		n, err := gy.port.Read(gy.buf[a:])
		if err != nil {
			return 0, 0, 0, err
		}
		a += n
	}

	var data []byte
	for i, b := range gy.buf {
		if b == datahead && i+datalen < bufsize {
			data = gy.buf[i : i+datalen]
			break
		}
	}
	if len(data) != datalen {
		return 0, 0, 0, errors.New("unexpected data length")
	}

	if data[0] != datahead && data[7] != datatail {
		return 0, 0, 0, errors.New("invalid data")
	}

	y := (int16(data[1]) << 8) | int16(data[2])
	p := (int16(data[3]) << 8) | int16(data[4])
	r := (int16(data[5]) << 8) | int16(data[6])
	return float64(y) / 100, float64(p) / 100, float64(r) / 100, nil
}

// Close ...
func (gy *GY25) Close() error {
	return gy.port.Close()
}
