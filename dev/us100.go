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
	"time"

	"github.com/stianeikeland/go-rpio"
	"github.com/tarm/serial"
)

type (
	// US100ModeType ...
	US100ModeType int

	// US100Option ...
	US100Option func(u *US100)
)

const (
	// US100UartMode ...
	US100UartMode US100ModeType = 1
	// US100TTLMode ...
	US100TTLMode US100ModeType = 2
)

var (
	trigData = []byte{0x55}
)

// US100 ...
type US100 struct {
	mode  US100ModeType
	buf   [4]byte
	retry int

	// ttl mode
	trig rpio.Pin
	echo rpio.Pin

	// uart mode
	name string
	baud int
	port *serial.Port
}

// US100Mode ...
func US100Mode(mode US100ModeType) US100Option {
	return func(u *US100) {
		u.mode = mode
	}
}

// US100TrigPin ...
func US100TrigPin(trig int8) US100Option {
	return func(u *US100) {
		u.trig = rpio.Pin(trig)
	}
}

// US100EchoPin ...
func US100EchoPin(echo int8) US100Option {
	return func(u *US100) {
		u.echo = rpio.Pin(echo)
	}
}

// US100Name ...
func US100Name(name string) US100Option {
	return func(u *US100) {
		u.name = name
	}
}

// US100Baud ...
func US100Baud(baud int) US100Option {
	return func(u *US100) {
		u.baud = baud
	}
}

// NewUS100 ...
func NewUS100(opts ...US100Option) *US100 {
	u := &US100{
		retry: 10,
	}
	for _, opt := range opts {
		opt(u)
	}

	if u.mode == US100TTLMode {
		u.trig.Output()
		u.trig.Low()
		u.echo.Input()
		return u
	}
	if u.mode == US100UartMode {
		if err := u.open(); err != nil {
			return nil
		}
		return u
	}
	return nil
}

// Dist is to measure the distance in cm
func (u *US100) Dist() float64 {
	if u.mode == US100UartMode {
		return u.DistByUart()
	}
	return u.DistByTTL()
}

// DistByUart is to measure the distance in cm
func (u *US100) DistByUart() float64 {
	if err := u.port.Flush(); err != nil {
		log.Printf("[us100]failed to flush serial, error: %v", err)
		return -1
	}
	// send trigger data
	n, err := u.port.Write(trigData)
	if n != 1 || err != nil {
		return -1
	}

	// read data
	n, err = u.port.Read(u.buf[:])
	if err != nil {
		u.Close()
		if err := u.open(); err != nil {
			log.Printf("[us100]failed to open serial, error: %v", err)
		}
		return -1
	}
	// check data len
	if n != 2 {
		log.Printf("[us100]incorrect data len, len: %v, expected: 2", n)
		return -1
	}
	// calc distance in cm
	return float64((uint16(u.buf[0])<<8)|uint16(u.buf[1])) / 10.0
}

// DistByTTL is to measure the distance in cm
func (u *US100) DistByTTL() float64 {
	u.trig.Low()
	u.delay(1)
	u.trig.High()
	u.delay(5)

	u.echo.PullDown()
	u.echo.Detect(rpio.RiseEdge)
	for !u.echo.EdgeDetected() {
		u.delay(1)
	}

	start := time.Now()
	u.echo.Detect(rpio.FallEdge)
	for !u.echo.EdgeDetected() {
		u.delay(1)
	}
	dist := time.Now().Sub(start).Seconds() * voiceSpeed / 2.0
	u.echo.Detect(rpio.NoEdge)
	u.trig.Low()
	return dist
}

// Close ...
func (u *US100) Close() {
	if u.mode == US100UartMode {
		u.port.Close()
	}
}

func (u *US100) open() error {
	c := &serial.Config{Name: u.name, Baud: u.baud}
	port, err := serial.OpenPort(c)
	if err != nil {
		return err
	}
	u.port = port
	return nil
}

// delay is to delay us microsecond
func (u *US100) delay(us int) {
	time.Sleep(time.Duration(us) * time.Microsecond)
}
