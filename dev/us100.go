/*
US-100 is an ultrasonic distance meter used to measure the distance to objects.
US-100 works in both modes of UART and Electrical Level(TTL).
TTL mode is used by default if you don't specify a mode for it.

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
 - ...............................................
 - !!! NOTE: TX->TXD, RX-RXD, NOT TX->RXD, RX-TXD
 - ...............................................
 - Trig/TX: must connect to GPIO-14 (TXD)
 - Echo/RX: must connect to GPIO-15 (RXD)

*/
package dev

import (
	"fmt"
	"time"

	"github.com/stianeikeland/go-rpio/v4"
	"github.com/tarm/serial"
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

// NewUS100TTL creates US100 using TTL mode
func NewUS100TTL(trig, echo uint8) (*US100, error) {
	us := &US100{
		mode: TTLMode,
		trig: rpio.Pin(trig),
		echo: rpio.Pin(echo),
	}
	us.trig.Output()
	us.trig.Low()
	us.echo.Input()
	return us, nil
}

// NewUS100Uart creates US100 using Uart mode
func NewUS100Uart(dev string, baud int) (*US100, error) {
	cfg := &serial.Config{
		Name:        dev,
		Baud:        baud,
		ReadTimeout: 1 * time.Second,
	}
	port, err := serial.OpenPort(cfg)
	if err != nil {
		return nil, err
	}
	return &US100{
		mode: UartMode,
		port: port,
	}, nil
}

// Value returns the distance in cm to objects
func (us *US100) Dist() (float64, error) {
	if us.mode == UartMode {
		return us.distFromUart()
	}
	return us.distFromTTL()
}

func (us *US100) distFromUart() (float64, error) {
	if err := us.port.Flush(); err != nil {
		return 0, fmt.Errorf("flush port error: %w", err)
	}
	// send trigger data
	n, err := us.port.Write(trigData)
	if n != 1 || err != nil {
		return 0, err
	}

	// read data
	a := 0
	for a < 2 {
		n, err := us.port.Read(us.buf[a:])
		if err != nil {
			return 0, fmt.Errorf("read port error: %w", err)
		}
		a += n
	}

	// check data len
	if a != 2 {
		return 0, fmt.Errorf("incorrect data len, len: %v, expected: 2", a)
	}
	// calc distance in cm
	return float64((uint16(us.buf[0])<<8)|uint16(us.buf[1])) / 10.0, nil
}

func (us *US100) distFromTTL() (float64, error) {
	us.trig.Low()
	delayUs(1)
	us.trig.High()
	delayUs(5)

	us.echo.PullDown()
	us.echo.Detect(rpio.RiseEdge)
	for !us.echo.EdgeDetected() {
		delayUs(1)
	}

	start := time.Now()
	us.echo.Detect(rpio.FallEdge)
	for !us.echo.EdgeDetected() {
		delayUs(1)
	}
	dist := time.Since(start).Seconds() * voiceSpeed / 2.0
	us.echo.Detect(rpio.NoEdge)
	us.trig.Low()
	return dist, nil
}

// Close ...
func (us *US100) Close() error {
	if us.mode == UartMode {
		return us.port.Close()
	}
	return nil
}
