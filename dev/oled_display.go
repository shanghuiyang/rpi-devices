/*
OledDisplay is an oled display module used to display text/image driving by ssd1306 driver.

Config Raspberry Pi:
1. $ sudo apt-get install -y python-smbus
2. $ sudo apt-get install -y i2c-tools
3. $ sudo raspi-config
4. 	-> [5 interface options] -> [P5 I2C] -> [yes] -> [ok]
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
	fontFile = "casio-fx-9860gii.ttf"
)

// OledDisplay ...
type OledDisplay struct {
	oled   *monochromeoled.OLED
	width  int
	height int
}

// NewOledDisplay ...
func NewOledDisplay(width, heigth int) (*OledDisplay, error) {
	oled, err := monochromeoled.Open(&i2c.Devfs{Dev: oledDev}, oledAddr, width, heigth)
	if err != nil {
		return nil, err
	}
	return &OledDisplay{
		oled:   oled,
		width:  width,
		height: heigth,
	}, nil
}

// DisplayImage displays an image on the screen
func (oled *OledDisplay) DisplayImage(img image.Image) error {
	if err := oled.oled.SetImage(0, 0, img); err != nil {
		return err
	}
	if err := oled.oled.Draw(); err != nil {
		return err
	}
	return nil
}

// Display displays the text on the screen.
// NOTE: It isn't implemented. It is here just for implementing the Display interface.
// Please draw your text to an image first, and then use DisplayImage()
func (oled *OledDisplay) DisplayText(text string, x, y int) error {
	return errors.New("not implement")
}

// On ...
func (oled *OledDisplay) On() error {
	return oled.oled.On()
}

// Off ...
func (oled *OledDisplay) Off() error {
	return oled.oled.Off()
}

// Clear ...
func (oled *OledDisplay) Clear() error {
	return oled.oled.Clear()
}

// Close ...
func (oled *OledDisplay) Close() error {
	_ = oled.oled.Clear()
	return oled.oled.Close()
}
