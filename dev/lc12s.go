/*
LC12S is wireless transceiver used to send and receive data via electromagnetic wave.
More details please ref to: https://world.taobao.com/item/594554513623.htm

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
 - VCC: any 3.3v
 - GND: any gnd pin
 - CS: any data pin. high-level: sleep, low-level: work
 - RX: must connect to GPIO-14/TXD
 - TX: must connect to GPIO-15/RXD

*/
package dev

import (
	"fmt"
	"io"
	"time"

	"github.com/stianeikeland/go-rpio/v4"
	"github.com/tarm/serial"
)

const (
	bufsz = 128
)

// LC12S implement Wireless interface
type LC12S struct {
	csPin rpio.Pin
	port  *serial.Port
}

// NewLC12S ...
func NewLC12S(dev string, baud int, csPin uint8) (*LC12S, error) {
	cfg := &serial.Config{
		Name:        dev,
		Baud:        baud,
		ReadTimeout: 3 * time.Second,
	}
	port, err := serial.OpenPort(cfg)
	if err != nil {
		return nil, err
	}
	l := &LC12S{
		csPin: rpio.Pin(csPin),
		port:  port,
	}
	l.csPin.Output()
	l.Sleep()
	return l, nil
}

// Send ...
func (l *LC12S) Send(data []byte) error {
	n, err := l.port.Write(data)
	if err != nil {
		return fmt.Errorf("write port error: %w", err)
	}

	if n != len(data) {
		return fmt.Errorf("not all data was sent")
	}

	return nil
}

// Receive ...
func (l *LC12S) Receive() ([]byte, error) {
	if err := l.port.Flush(); err != nil {
		return nil, fmt.Errorf("flush port error: %w", err)
	}

	var buf [bufsz]byte
	n, err := l.port.Read(buf[:])
	if err == io.EOF {
		return []byte{}, nil
	}
	if err != nil {
		return nil, err
	}
	if n > bufsz {
		return nil, fmt.Errorf("buf overflow")
	}
	return buf[:n], nil
}

// Sleep ...
func (l *LC12S) Sleep() {
	l.csPin.High()
}

// Wakeup ...
func (l *LC12S) Wakeup() {
	l.csPin.Low()
}

// Close ...
func (l *LC12S) Close() error {
	return l.port.Close()
}
