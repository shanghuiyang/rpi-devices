/*
US-100 is an ultrasonic distance meter used to measure the distance to objects.
US-100 supports both of interfaces: GPIO and UART.
min distance: 2cm
max distance: 450cm

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
GPIO Interface:
  - VCC: any 3.3v or 5v pin
  - GND: any gnd pin
  - Trig: any gnd pin
  - Echo: any gnd pin

    UART Interface:
  - VCC: any 3.3v or 5v pin
  - GND: any gnd pin
  - ...............................................
  - !!! NOTE: TX->TXD, RX-RXD, NOT TX->RXD, RX-TXD
  - ...............................................
  - TX: must connect to GPIO-14 (TXD)
  - RX: must connect to GPIO-15 (RXD)

*/
package dev

import (
	"errors"
	"fmt"
	"sort"
	"time"

	"github.com/stianeikeland/go-rpio/v4"
	"github.com/tarm/serial"
)

const (
	us100MaxDistance = 699  // cm
	us100Timeout     = 1000 // Nanosecond
	us100MaxRetries  = 5
	us100Trim        = 1
)

var (
	trigData = []byte{0x55}
)

// US100 ...
type US100 struct {
	iface InterfaceType
	buf   [4]byte

	// ttl mode
	trig rpio.Pin
	echo rpio.Pin

	// uart mode
	port *serial.Port
}

// NewUS100GPIO creates US100 using GPOI interface
func NewUS100GPIO(trig, echo uint8) (*US100, error) {
	us := &US100{
		iface: GPIO,
		trig:  rpio.Pin(trig),
		echo:  rpio.Pin(echo),
	}
	us.trig.Output()
	us.trig.Low()
	us.echo.Input()
	return us, nil
}

// NewUS100UART creates US100 using UART interface
func NewUS100UART(dev string, baud int) (*US100, error) {
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
		iface: UART,
		port:  port,
	}, nil
}

// Value returns the distance in cm to objects
func (us *US100) Dist() (float64, error) {
	if us.iface == UART {
		return us.distFromUART()
	}

	values := make([]float64, 0, hcsr04MaxRetries)
	for i := 0; i < us100MaxRetries; i++ {
		dist, err := us.distFromGPIO()
		if err != nil {
			time.Sleep(1 * time.Microsecond)
			continue
		}
		values = append(values, dist)
		time.Sleep(1 * time.Microsecond)
	}

	if len(values) <= us100Trim*2 {
		return 0, errors.New("device not ready")
	}

	sort.Float64s(values)
	trimmed := values[us100Trim : len(values)-us100Trim]
	var sum float64
	for _, v := range trimmed {
		sum += v
	}
	return sum / float64(len(trimmed)), nil
}

func (us *US100) distFromUART() (float64, error) {
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

func (us *US100) distFromGPIO() (float64, error) {
	us.trig.Low()
	delayUs(1)
	us.trig.High()
	delayUs(1)

	us.echo.PullDown()
	us.echo.Detect(rpio.RiseEdge)
	for i := 0; !us.echo.EdgeDetected(); i++ {
		if i >= us100Timeout {
			return us100MaxDistance, nil
		}
		delayNs(1)
	}
	start := time.Now()
	us.echo.Detect(rpio.FallEdge)
	for i := 0; !us.echo.EdgeDetected(); i++ {
		if i >= us100Timeout {
			return us100MaxDistance, nil
		}
		delayNs(1)
	}
	dist := time.Since(start).Seconds() * voiceSpeed / 2.0
	us.echo.Detect(rpio.NoEdge)
	us.trig.Low()
	return dist, nil
}

// Close ...
func (us *US100) Close() error {
	if us.iface == UART {
		return us.port.Close()
	}
	return nil
}
