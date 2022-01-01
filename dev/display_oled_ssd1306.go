/*
SSD1306Display is a driver for the oled display module drived by SSD1306 chip.
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
	30: -- -- -- -- -- -- -- -- -- -- -- -- 3c -- -- --
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

	"github.com/mdp/monochromeoled"
	"golang.org/x/exp/io/i2c"
)

const (
	oledDev  = "/dev/i2c-1"
	oledAddr = 0x3c
)

// SSD1306Display is a driver for the oled display module drived by SSD1306 chip.
// It is an implement of Display interface.
type SSD1306Display struct {
	oled   *monochromeoled.OLED
	width  int
	height int
}

// NewSSD1306Display creates a driver for the oled display module drived by SSD1306 chip
func NewSSD1306Display(width, heigth int) (*SSD1306Display, error) {
	oled, err := monochromeoled.Open(&i2c.Devfs{Dev: oledDev}, oledAddr, width, heigth)
	if err != nil {
		return nil, err
	}
	return &SSD1306Display{
		oled:   oled,
		width:  width,
		height: heigth,
	}, nil
}

// Image displays an image on the screen
func (display *SSD1306Display) Image(img image.Image) error {
	if err := display.oled.SetImage(0, 0, img); err != nil {
		return err
	}
	if err := display.oled.Draw(); err != nil {
		return err
	}
	return nil
}

// Text displays the text on the screen.
// NOTE: It isn't implemented. It is here just for implementing the Display interface.
// Please draw your text to an image first, and then use DisplayImage()
func (display *SSD1306Display) Text(text string, x, y int) error {
	return errors.New("not implement")
}

// On ...
func (display *SSD1306Display) On() error {
	return display.oled.On()
}

// Off ...
func (display *SSD1306Display) Off() error {
	return display.oled.Off()
}

// Clear ...
func (display *SSD1306Display) Clear() error {
	return display.oled.Clear()
}

// Close ...
func (display *SSD1306Display) Close() error {
	_ = display.oled.Clear()
	return display.oled.Close()
}
