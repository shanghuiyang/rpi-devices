/*
Package dev ...

US-100 is an ultrasonic distance meter,
which can measure the distance to the an object like a box.
US-100 works in both modes of UART and Electrical Level.
This program uses UART mode to drive the module.

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
 - VCC: any 3.3v or 5v pin
 - GND: any gnd pin
 - ...............................................
 - !!! NOTE: TX->TXD, RX-RXD, NOT TX->RXD, RX-TXD
 - ...............................................
 - Trig/TX: must connect to pin  8(gpio 14) (TXD)
 - Echo/RX: must connect to pin 10(gpio 15) (RXD)
*/
package dev

import (
	"log"

	"github.com/tarm/serial"
)

const (
	logTagUS100 = "US100"
)

var (
	trigData = []byte{0x55}
)

// US100 ...
type US100 struct {
	port  *serial.Port
	buf   [32]byte
	retry int
}

// NewUS100 ...
func NewUS100() *US100 {
	u := &US100{
		retry: 10,
	}
	if err := u.open(); err != nil {
		return nil
	}
	return u
}

// Dist is to measure the distance in cm
func (u *US100) Dist() float64 {
	if err := u.port.Flush(); err != nil {
		log.Printf("failed to flush serial, error: %v", err)
		return -1
	}
	// send trigger data
	n, err := u.port.Write(trigData)
	if n != 1 || err != nil {
		return -1
	}

	// read data
	p := 0
	for p < 2 {
		n, err := u.port.Read(u.buf[p:])
		if err != nil {
			u.Close()
			if err := u.open(); err != nil {
				log.Printf("failed to open serial, error: %v", err)
			}
			return -1
		}
		p += n
		log.Printf("[us100]read data, p=%v", p)
	}
	// check data len
	if p != 2 {
		log.Printf("incorrect data len, len: %v, expected: 2", p)
		return -1
	}
	// calc distance in cm
	return float64((uint16(u.buf[0])<<8)|uint16(u.buf[1])) / 10.0
}

// Close ...
func (u *US100) Close() {
	u.port.Close()
}

func (u *US100) open() error {
	c := &serial.Config{Name: "/dev/ttyAMA0", Baud: 9600}
	port, err := serial.OpenPort(c)
	if err != nil {
		return err
	}
	u.port = port
	return nil
}
