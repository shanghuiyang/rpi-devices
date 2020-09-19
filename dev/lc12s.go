/*
Package dev ...

LC12S is wireless transceiver module.
https://world.taobao.com/item/594554513623.htm

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
 - VCC: any 3.3v
 - GND: any gnd pin
 - CS: any data pin. high-level: sleep, low-level: work
 - TX: must connect to pin 10(gpio 15) (RXD)
 - RX: must connect to pin  8(gpio 14) (TXD)
*/
package dev

import (
	"fmt"
	"io"
	"log"

	"github.com/stianeikeland/go-rpio"
	"github.com/tarm/serial"
)

const (
	bufsz = 128
)

// LC12S ...
type LC12S struct {
	csPin rpio.Pin
	port  *serial.Port
}

// NewLC12S ...
func NewLC12S(csPin uint8) *LC12S {
	l := &LC12S{
		csPin: rpio.Pin(csPin),
	}
	if err := l.open(); err != nil {
		return nil
	}
	l.csPin.Output()
	l.Sleep()
	return l
}

// Write ...
func (l *LC12S) Write(data []byte) error {
	n, err := l.port.Write(data)
	if err != nil {
		return err
	}

	if n != len(data) {
		return fmt.Errorf("[lc12s]wrote % btyes data, but expect % bytes", n, len(data))
	}

	return nil
}

// Read ...
func (l *LC12S) Read() ([]byte, error) {
	if err := l.port.Flush(); err != nil {
		log.Printf("[lc12s]failed to flush serial, error: %v", err)
		return nil, err
	}

	var buf [bufsz]byte
	n, err := l.port.Read(buf[:])
	if err == io.EOF {
		return []byte{}, nil
	}
	if err != nil {
		l.Close()
		// re-open
		if err := l.open(); err != nil {
			log.Printf("[lc12s]failed to open serial, error: %v", err)
		}
		return nil, err
	}
	if n > bufsz {
		return nil, fmt.Errorf("[lc12s]received % bytes data, expect less than %v bytes", n, bufsz)
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
func (l *LC12S) Close() {
	l.port.Close()
}

func (l *LC12S) open() error {
	c := &serial.Config{
		Name: "/dev/ttyAMA0",
		Baud: 9600,
	}
	port, err := serial.OpenPort(c)
	if err != nil {
		return err
	}
	l.port = port
	return nil
}
