/*
IRCoder is an infrared(IR) encoder and decoder

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
 - VCC: any 5v pin
 - GND: any gnd pin
 - RXT: must connect to GPIO-14/TXD
 - TXD: must connect to GPIO-15/RXD

*/
package dev

import (
	"fmt"
	"time"

	"github.com/tarm/serial"
)

// IRCoder ...
type IRCoder struct {
	port *serial.Port
}

// NewIRCoder ...
func NewIRCoder(dev string, baud int) *IRCoder {
	ir := &IRCoder{}
	if err := ir.open(dev, baud); err != nil {
		return nil
	}
	return ir
}

func (ir *IRCoder) Send(data []byte) error {
	if err := ir.port.Flush(); err != nil {
		return fmt.Errorf("port flush error: %w", err)
	}

	n, err := ir.port.Write(data)
	if err != nil {
		return fmt.Errorf("port write error: %w", err)
	}

	if n != len(data) {
		return fmt.Errorf("not all data was sent")
	}

	return nil
}

// Close ...
func (ir *IRCoder) Close() error {
	return ir.port.Close()
}

func (ir *IRCoder) open(dev string, baud int) error {
	c := &serial.Config{
		Name:        dev,
		Baud:        baud,
		ReadTimeout: 3 * time.Second,
	}
	port, err := serial.OpenPort(c)
	if err != nil {
		return err
	}
	ir.port = port
	return nil
}
