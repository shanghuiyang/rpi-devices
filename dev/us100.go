/*
US-100 is an ultrasonic distance meter used to measure the distance to objects.
US-100 works in both modes of UART and Electrical Level(TTL).
TTL mode is used by default if you don't specify a mode for it.

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
	"fmt"
	"log"
	"time"

	"github.com/stianeikeland/go-rpio/v4"
	"github.com/tarm/serial"
)

const (
	defaultUS100Name = "/dev/ttyAMA0"
	defaultUS100Baud = 9600
)

var (
	trigData = []byte{0x55}
)

// US100Config ...
type US100Config struct {
	Mode  ComMode
	Trig  int8
	Echo  int8
	Dev   string
	Baud  int
	Retry int
}

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

	// UART mode
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

// Value returns the distance in cm to objects
func (u *US100) Dist() (float64, error) {
	if u.mode == UartMode {
		return u.distByUart()
	}
	return u.distByTTL()
}

func (u *US100) distByUart() (float64, error) {
	if err := u.port.Flush(); err != nil {
		log.Printf("[us100]failed to flush serial, error: %v", err)
		return 0, err
	}
	// send trigger data
	n, err := u.port.Write(trigData)
	if n != 1 || err != nil {
		return 0, err
	}

	// read data
	a := 0
	for a < 2 {
		n, err := u.port.Read(u.buf[a:])
		if err != nil {
			log.Printf("[us100]failed to read serial, error: %v", err)
			return 0, err
		}
		a += n
	}

	// check data len
	if a != 2 {
		log.Printf("[us100]incorrect data len, len: %v, expected: 2", a)
		return 0, fmt.Errorf("incorrect data len, len: %v, expected: 2", a)
	}
	// calc distance in cm
	return float64((uint16(u.buf[0])<<8)|uint16(u.buf[1])) / 10.0, nil
}

func (u *US100) distByTTL() (float64, error) {
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
	dist := time.Since(start).Seconds() * voiceSpeed / 2.0
	u.echo.Detect(rpio.NoEdge)
	u.trig.Low()
	return dist, nil
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
