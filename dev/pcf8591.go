/*
PCF8591 is the driver of PCF8591 module.

Jumper:
- remove jumpers on P4 & P5, keep the jumper on P6

Config Raspberry Pi:
1. $ sudo apt-get install -y python-smbus
2. $ sudo apt-get install -y i2c-tools
3. $ sudo raspi-config
4. 	-> [5 Interface Options] -> [P5 I2C] -> [yes] -> [ok]
5. $ sudo reboot now
6. check: $ sudo i2cdetect -y 1
	it works if you saw following message:
	~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
	     0  1  2  3  4  5  6  7  8  9  a  b  c  d  e  f
	00:          -- -- -- -- -- -- -- -- -- -- -- -- --
	10: -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- --
	20: -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- --
	30: -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- --
	40: -- -- -- -- -- -- -- -- 48 -- -- -- -- -- -- --
	50: -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- --
	60: -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- --
	70: -- -- -- -- -- -- -- --
	~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

Connect to Raspberry Pi:
 - VCC: any 3.3v pin
 - GND: any GND pin
 - SDA: GPIO-2 (SDA)
 - SCL: GPIO-3 (SCL)

*/
package dev

import (
	"fmt"

	"golang.org/x/exp/io/i2c"
)

const (
	pcf8591Dev  = "/dev/i2c-1"
	pcf8591Addr = 0x48
	ctrAIN0     = 0x40
	ctrAIN1     = 0x41
	ctrAIN2     = 0x42
	ctrAIN3     = 0x43
)

// PCF8591 ...
type PCF8591 struct {
	dev *i2c.Device
}

// NewPCF8591 ...
func NewPCF8591() (*PCF8591, error) {
	dev, err := i2c.Open(&i2c.Devfs{Dev: pcf8591Dev}, pcf8591Addr)
	if err != nil {
		return nil, err
	}
	return &PCF8591{
		dev: dev,
	}, nil
}

// ReadAIN0 ...
func (pcf *PCF8591) ReadAIN0() ([]byte, error) {
	if err := pcf.dev.Write([]byte{ctrAIN0}); err != nil {
		return nil, fmt.Errorf("i2c write error: %w", err)
	}
	data := make([]byte, 1)
	if err := pcf.dev.Read(data); err != nil {
		return nil, fmt.Errorf("i2c read error: %w", err)
	}
	return data, nil
}

// ReadAIN1 ...
func (pcf *PCF8591) ReadAIN1() ([]byte, error) {
	if err := pcf.dev.Write([]byte{ctrAIN1}); err != nil {
		return nil, fmt.Errorf("i2c write error: %w", err)
	}
	data := make([]byte, 1)
	if err := pcf.dev.Read(data); err != nil {
		return nil, fmt.Errorf("i2c read error: %w", err)
	}
	return data, nil
}

// ReadAIN2 ...
func (pcf *PCF8591) ReadAIN2() ([]byte, error) {
	if err := pcf.dev.Write([]byte{ctrAIN2}); err != nil {
		return nil, fmt.Errorf("i2c write error: %w", err)
	}
	data := make([]byte, 1)
	if err := pcf.dev.Read(data); err != nil {
		return nil, fmt.Errorf("i2c read error: %w", err)
	}
	return data, nil
}

// ReadAIN3 ...
func (pcf *PCF8591) ReadAIN3() ([]byte, error) {
	if err := pcf.dev.Write([]byte{ctrAIN3}); err != nil {
		return nil, fmt.Errorf("i2c write error: %w", err)
	}
	data := make([]byte, 1)
	if err := pcf.dev.Read(data); err != nil {
		return nil, fmt.Errorf("i2c read error: %w", err)
	}
	return data, nil
}

// Close ...
func (pcf *PCF8591) Close() error {
	return pcf.dev.Close()
}
