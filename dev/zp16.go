/*
ZP16 is a gas detect module.

Config Raspberry Pi:
1. $ sudo raspi-config
	-> [P5 interface] -> P6 Serial: disable -> [no] -> [yes]
2. $ sudo vim /boot/config.txt
	add following two lines:
	~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
	enable_uart=1
	dtoverlay=pi3-miniuart-bt
	~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
3. $ sudo reboot now
4. $ sudo cat /dev/ttyAMA0
	should see somethings output

Connect to Raspberry Pi:
 - VCC: any 3.3v or 5v pin
 - GND: any gnd pin
 - RXT: must connect to GPIO-14/TXD
 - TXD: must connect to GPIO-15/RXD
*/
package dev

import (
	"fmt"

	"github.com/tarm/serial"
)

// ZP16 ...
type ZP16 struct {
	port *serial.Port
}

// NewZP16 ...
func NewZP16(dev string, baud int) (*ZP16, error) {
	cfg := &serial.Config{
		Name: dev,
		Baud: baud,
	}
	port, err := serial.OpenPort(cfg)
	if err != nil {
		return nil, err
	}
	return &ZP16{port}, nil
}

// Loc ...
func (zp *ZP16) CO() (float64, error) {
	if err := zp.port.Flush(); err != nil {
		return 0, fmt.Errorf("flush port error: %w", err)
	}

	// read data
	a := 0
	var buf [16]byte
	// log.Printf("start to read")
	// bytes:	| byte0 	byte1 	byte2 	byte3 	byte4 	byte5 	byte6 		byte7 		byte8
	// ---------+------------------------------------------------------------------------------------
	// example:	|  ff		 34 	 11 	 02 	 03      e8 	 03 		 e8 		 e3
	// desc:	| start		name	unit	point	high	low		full-high	full-low	checksum
	for a < 9 {
		n, err := zp.port.Read(buf[a:])
		if err != nil {
			return 0, fmt.Errorf("read port error: %w", err)
		}
		// log.Printf("n=%v", n)
		a += n
	}
	// log.Printf("read done")
	// for _, b :=range buf {
	// 	log.Printf("%x", b)
	// }
	if buf[0] != 0xff {
		return 0, fmt.Errorf("invalid data")
	}
	// check data len
	if a != 9 {
		return 0, fmt.Errorf("incorrect data len")
	}

	// var checksum uint16
	// for i := 0; i < 8; i++ {
	// 	checksum += uint16(buf[i])
	// }
	// if (^checksum)+1 != uint16(buf[8]) {
	// 	return 0, fmt.Errorf("checksum failed")
	// }
	// log.Printf("b3=%d, b4=%d, b5=%v", buf[3], buf[4], buf[5])
	co := float64((uint16(buf[4]) << 8) | uint16(buf[5]))
	for i := 0; i < int(buf[3]); i++ {
		co /= 10.0
	}
	return co, nil
}

// Close ...
func (zp *ZP16) Close() error {
	return zp.port.Close()
}
