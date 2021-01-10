/*
Package dev ...

GY25 is the driver of GY25, an angle sensor which can be used to detect yaw, pitch and roll angle.

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

	"github.com/tarm/serial"
)

const (
	datalen  = 8
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

// GY25 ...
type GY25 struct {
	port *serial.Port
	buf  [16]byte
}

// NewGY25 ...
func NewGY25() *GY25 {
	g := &GY25{}
	if err := g.open(); err != nil {
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

	n, err := g.port.Read(g.buf[:])
	if err != nil {
		return 0, 0, 0, err
	}

	if n != datalen {
		return 0, 0, 0, fmt.Errorf("incorrect data len: %v, expected %v", n, datalen)
	}

	if g.buf[0] != datahead && g.buf[7] != datatail {
		return 0, 0, 0, fmt.Errorf("invalid data")
	}

	yaw := (int16(g.buf[1]) << 8) | int16(g.buf[2])
	pitch := (int16(g.buf[3]) << 8) | int16(g.buf[4])
	roll := (int16(g.buf[5]) << 8) | int16(g.buf[6])
	return float64(yaw) / 100, float64(pitch) / 100, float64(roll) / 100, nil
}

// Close ...
func (g *GY25) Close() {
	g.port.Close()
}

func (g *GY25) open() error {
	c := &serial.Config{
		Name: "/dev/ttyAMA0",
		Baud: 115200,
	}
	port, err := serial.OpenPort(c)
	if err != nil {
		return err
	}
	g.port = port
	return nil
}
