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
	"log"
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
		log.Printf("failed to flush serial, error: %v", err)
		return err
	}

	n, err := ir.port.Write(data)
	if n != 5 || err != nil {
		return err
	}

	if n != 5 {
		return fmt.Errorf("send %v bytes data, expected 5 bytes", n)
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
