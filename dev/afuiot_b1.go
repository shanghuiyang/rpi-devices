/*
Afuiot-B1 is a sensor used to measure the temperature of human or object through infrared ray.

Config Raspberry Pi:
1. $ sudo raspi-config
	-> [P5 interface] -> P6 Serial: disable -> [no] -> [yes]
2. $ sudo vim /boot/config.txt
	add following two lines:
	~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
	enable_uart=1
	dtoverlay=pi3-miniuart-bt
	~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
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
	"fmt"

	"github.com/tarm/serial"
)

var (
	SetAfuiotB1Uart       = []byte{0xAA, 0xA5, 0x04, 0x05, 0x01, 0x0A, 0x55}
	SetAfuiotB1I2C        = []byte{0xAA, 0xA5, 0x04, 0x05, 0x02, 0x0A, 0x55}
	SetAfuiotB1ObjectTemp = []byte{0xAA, 0xA5, 0x04, 0x02, 0x01, 0x07, 0x55}
	SetAfuiotB1HumanTemp  = []byte{0xAA, 0xA5, 0x04, 0x02, 0x02, 0x08, 0x55}

	cmdAfuiotB1Measure = []byte{0xAA, 0xA5, 0x03, 0x01, 0x04, 0x55}
)

// AfuiotB1 implements Thermometer interface
type AfuiotB1 struct {
	port *serial.Port
	buf  [8]byte
}

// NewAfuiotB1 ...
func NewAfuiotB1(dev string, baud int) (*AfuiotB1, error) {
	cfg := &serial.Config{
		Name: dev,
		Baud: baud,
	}
	port, err := serial.OpenPort(cfg)
	if err != nil {
		return nil, err
	}
	return &AfuiotB1{port: port}, nil
}

// Temperature ...
func (b1 *AfuiotB1) Temperature() (float64, error) {
	if err := b1.port.Flush(); err != nil {
		return 0, fmt.Errorf("flush port error: %w", err)
	}

	if err := b1.Set(cmdAfuiotB1Measure); err != nil {
		return 0, nil
	}

	n, err := b1.port.Read(b1.buf[:])
	if n != 8 {
		return 0, fmt.Errorf("read %v bytes, expected 8 bytes", n)
	}

	if err != nil {
		return 0, err
	}

	return float64((uint16(b1.buf[5])<<8)|uint16(b1.buf[6])) / 10.0, nil
}

// Set ...
func (b1 *AfuiotB1) Set(bytes []byte) error {
	n, err := b1.port.Write(bytes)
	if n != len(bytes) {
		return fmt.Errorf("write %v bytes, but expect %v bytes", n, len(bytes))
	}
	if err != nil {
		return err
	}

	return nil
}

// Close ...
func (b1 *AfuiotB1) Close() error {
	return b1.port.Close()
}
