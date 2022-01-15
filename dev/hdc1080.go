/*
HDC1080 is a sensor used to measure temperature and humidity. It is a implement of TempHumidity interface.

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
	40: 40 -- -- -- -- -- -- -- -- -- -- -- -- -- -- --
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
	"errors"

	"golang.org/x/exp/io/i2c"
)

const (
	hdc080Dev   = "/dev/i2c-1"
	hdc1080Addr = 0x40
	maxRetry    = 10
)

var (
	hdc1080Cmd = []byte{0x00}
)

// HDC1080 ...
type HDC1080 struct {
	dev *i2c.Device
}

// NewHDC1080 implement Thermohygrometer interface
func NewHDC1080() (*HDC1080, error) {
	dev, err := i2c.Open(&i2c.Devfs{Dev: hdc080Dev}, hdc1080Addr)
	if err != nil {
		return nil, err
	}
	return &HDC1080{dev: dev}, nil
}

// TempHumidity ...
func (hdc *HDC1080) TempHumidity() (temp, humi float64, err error) {
	if hdc.dev.Write(hdc1080Cmd); err != nil {
		return 0, 0, err
	}

	for i := 0; i < maxRetry; i++ {
		data := make([]byte, 4)
		if err := hdc.dev.Read(data); err != nil {
			delayMs(1)
			continue
		}
		t16 := uint16(data[0])<<8 | uint16(data[1])
		h16 := uint16(data[2])<<8 | uint16(data[3])
		t := 165*(float64(t16)/65536) - 40
		h := (float64(h16) / 65536) * 100
		return t, h, nil
	}
	return 0, 0, errors.New("HDC1080 isn't ready, please retry later")
}

// Close ...
func (hdc *HDC1080) Close() error {
	return hdc.dev.Close()
}
