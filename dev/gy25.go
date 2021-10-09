/*
GY25 is an accelerometer used to detect yaw, pitch and roll angles.

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
 - TXD: must connect to pin 10(gpio 15) (RXD)
 - RXT: must connect to pin  8(gpio 14) (TXD)

*/
package dev

import (
	"fmt"
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
func NewGY25(dev string, baud int) *GY25 {
	g := &GY25{}
	if err := g.open(dev, baud); err != nil {
		return nil
	}
	return g
}

// SetMode ...
func (g *GY25) SetMode(mode GY25Mode) error {
	if err := g.port.Flush(); err != nil {
		return err
	}

	n, err := g.port.Write(mode)
	if n != 2 || err != nil {
		return err
	}
	return nil
}

// Angles ...
func (g *GY25) Angles() (float64, float64, float64, error) {
	if err := g.port.Flush(); err != nil {
		return 0, 0, 0, err
	}

	a := 0
	for a < 16 {
		n, err := g.port.Read(g.buf[a:])
		if err != nil {
			return 0, 0, 0, err
		}
		a += n
	}

	var data []byte
	for i, b := range g.buf {
		if b == datahead && i+datalen < bufsize {
			data = g.buf[i : i+datalen]
			break
		}
	}
	if len(data) != datalen {
		return 0, 0, 0, fmt.Errorf("incorrect data len: %v, expected %v", len(data), datalen)
	}

	if data[0] != datahead && data[7] != datatail {
		return 0, 0, 0, fmt.Errorf("invalid data, validation failed")
	}

	yaw := (int16(data[1]) << 8) | int16(data[2])
	pitch := (int16(data[3]) << 8) | int16(data[4])
	roll := (int16(data[5]) << 8) | int16(data[6])
	return float64(yaw) / 100, float64(pitch) / 100, float64(roll) / 100, nil
}

// Close ...
func (g *GY25) Close() {
	g.port.Close()
}

func (g *GY25) open(dev string, baud int) error {
	c := &serial.Config{
		Name:        dev,
		Baud:        baud,
		ReadTimeout: 3 * time.Second,
	}
	port, err := serial.OpenPort(c)
	if err != nil {
		return err
	}
	g.port = port
	return nil
}
