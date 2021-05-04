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

const (
	defaultUS100Name = "/dev/ttyAMA0"
	defaultUS100Baud = 9600
)

var (
	trigData = []byte{0x55}
)

// US100 ...
type US100 struct {
	mode ComMode
	buf  [4]byte

	// ttl mode
	trig rpio.Pin
	echo rpio.Pin

	// uart mode
	port *serial.Port
}

// NewUS100 ...
func NewUS100(cfg *US100Config) *US100 {
	u := &US100{
		mode: cfg.Mode,
	}

	if u.mode == TTLMode {
		u.trig = rpio.Pin(cfg.Trig)
		u.echo = rpio.Pin(cfg.Echo)
		u.trig.Output()
		u.trig.Low()
		u.echo.Input()
		return u
	}
	if u.mode == UartMode {
		dev := cfg.Dev
		baud := cfg.Baud

		if cfg.Dev == "" {
			dev = defaultUS100Name
		}
		if cfg.Baud == 0 {
			baud = defaultUS100Baud
		}
		if err := u.open(dev, baud); err != nil {
			return nil
		}
		return u
	}
	return nil
}

// Dist is to measure the distance in cm
func (u *US100) Dist() float64 {
	if u.mode == UartMode {
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
	a := 0
	for a < 2 {
		n, err := u.port.Read(u.buf[a:])
		if err != nil {
			log.Printf("[us100]failed to read serial, error: %v", err)
			return -1
		}
		a += n
	}

	// check data len
	if a != 2 {
		log.Printf("[us100]incorrect data len, len: %v, expected: 2", a)
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
	if u.mode == UartMode {
		u.port.Close()
	}
}

func (u *US100) open(dev string, baud int) error {
	c := &serial.Config{Name: dev, Baud: baud}
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
