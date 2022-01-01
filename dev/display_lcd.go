/*
LcdDisplay is a driver for LCD Dispaly.
Please NOTE that I only test it on a 1602A lcd display module.

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
	20: -- -- -- -- -- -- -- 27 -- -- -- -- -- -- -- --
	30: -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- --
	40: -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- --
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
	"image"

	"golang.org/x/exp/io/i2c"
)

const (
	lcdDev               = "/dev/i2c-1"
	lcdAddr              = 0x27
	lcdEnable       byte = 0x04
	lcdBacklightOn  byte = 0x08
	lcdBacklightOff byte = 0x00
)

// LcdDisplay is a driver for LCD Dispaly.
type LcdDisplay struct {
	width  int
	height int
	blkOn  bool
	dev    *i2c.Device
}

// NewLcdDisplay creates a driver for LCD display.
// It is an implement of Display interface.
// Please NOTE that I only test it on a 1602A lcd display module.
func NewLcdDisplay(width, height int) (*LcdDisplay, error) {
	dev, err := i2c.Open(&i2c.Devfs{Dev: lcdDev}, lcdAddr)
	if err != nil {
		return nil, err
	}
	lcd := &LcdDisplay{
		width:  width,
		height: height,
		dev:    dev,
	}
	lcd.sendCommand(0x33) // initialise
	delayMs(1)

	lcd.sendCommand(0x32) // initialise
	delayMs(1)

	lcd.sendCommand(0x06) // cursor move direction
	delayMs(1)

	lcd.sendCommand(0x0C) // display on, cursor off, blink off
	delayMs(1)

	lcd.sendCommand(0x28) // data length, number of lines, font size
	delayMs(1)

	lcd.sendCommand(0x01) // Clear display
	delayMs(1)

	return lcd, nil
}

// Image displays an image on the screen.
// NOTE: it isn't be implemented yet.
func (lcd *LcdDisplay) Image(img image.Image) error {
	return errors.New("not implement")
}

// Text display text on the screen
func (lcd *LcdDisplay) Text(text string, x, y int) error {
	if x < 0 {
		x = 0
	}
	if x >= lcd.width {
		x = lcd.width - 1
	}

	if y < 0 {
		y = 0
	}
	if y >= lcd.height {
		y = lcd.height - 1
	}

	startPos := byte(0x80 + 0x40*y + x)
	if err := lcd.sendCommand(startPos); err != nil {
		return err
	}
	btext := []byte(text)
	for _, c := range btext {
		if err := lcd.sendData(c); err != nil {
			return err
		}
	}
	return nil
}

func (lcd *LcdDisplay) On() error {
	lcd.blkOn = true
	return lcd.dev.Write([]byte{lcdBacklightOn})
}

func (lcd *LcdDisplay) Off() error {
	lcd.blkOn = false
	return lcd.dev.Write([]byte{lcdBacklightOff})
}

func (lcd *LcdDisplay) Clear() error {
	return lcd.sendCommand(0x01)
}

func (lcd *LcdDisplay) Close() error {
	if err := lcd.Clear(); err != nil {
		return err
	}
	if err := lcd.Off(); err != nil {
		return err
	}
	if err := lcd.dev.Close(); err != nil {
		return err
	}
	return nil
}

func (lcd *LcdDisplay) sendCommand(b byte) error {
	backlight := lcdBacklightOff
	if lcd.blkOn {
		backlight = lcdBacklightOn
	}
	high := (b & 0xF0) | backlight
	low := ((b << 4) & 0xF0) | backlight

	if err := lcd.dev.Write([]byte{high}); err != nil {
		return err
	}
	if err := lcd.triggeEnable(high); err != nil {
		return err
	}

	if err := lcd.dev.Write([]byte{low}); err != nil {
		return err
	}
	if err := lcd.triggeEnable(low); err != nil {
		return err
	}
	return nil
}

func (lcd *LcdDisplay) sendData(b byte) error {
	backlight := lcdBacklightOff
	if lcd.blkOn {
		backlight = lcdBacklightOn
	}
	high := 0x1 | (b & 0xF0) | backlight
	low := 0x1 | ((b << 4) & 0xF0) | backlight

	if err := lcd.dev.Write([]byte{high}); err != nil {
		return err
	}
	if err := lcd.triggeEnable(high); err != nil {
		return err
	}

	if err := lcd.dev.Write([]byte{low}); err != nil {
		return err
	}
	if err := lcd.triggeEnable(low); err != nil {
		return err
	}
	return nil
}

func (lcd *LcdDisplay) triggeEnable(b byte) error {
	delayMs(1)
	buf := b | lcdEnable
	if err := lcd.dev.Write([]byte{buf}); err != nil {
		return err
	}

	delayMs(1)
	buf = b & ^lcdEnable
	if err := lcd.dev.Write([]byte{buf}); err != nil {
		return err
	}
	delayMs(1)
	return nil
}
