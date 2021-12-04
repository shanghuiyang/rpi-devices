/*
OledDisplay is an oled display module used to display text/image driving by ssd1306 driver.

Config Raspberry Pi:
1. $ sudo apt-get install -y python-smbus
2. $ sudo apt-get install -y i2c-tools
3. $ sudo raspi-config
4. 	-> [5 interface options] -> [p5 i2c] ->[yes] -> [ok]
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
	"image"
	"image/draw"
	"io/ioutil"

	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"github.com/mdp/monochromeoled"
	"golang.org/x/exp/io/i2c"
)

const (
	fontFile = "casio-fx-9860gii.ttf"
)

// OledDisplay ...
type OledDisplay struct {
	*monochromeoled.OLED
	width  int
	height int
	font   *truetype.Font
}

// NewOledDisplay ...
func NewOledDisplay(width, heigth int) (*OledDisplay, error) {
	oled, err := monochromeoled.Open(&i2c.Devfs{Dev: "/dev/i2c-1"}, 0x3c, width, heigth)
	if err != nil {
		return nil, err
	}
	a, err := ioutil.ReadFile(fontFile)
	if err != nil {
		return nil, err
	}
	font, err := truetype.Parse(a)
	if err != nil {
		return nil, err
	}
	return &OledDisplay{
		OLED:   oled,
		width:  width,
		height: heigth,
		font:   font,
	}, nil
}

// Display ...
func (oled *OledDisplay) Display(text string, fontSize float64, x, y int) error {
	image, err := oled.drawText(text, fontSize, x, y)
	if err != nil {
		return err
	}
	if err := oled.OLED.SetImage(0, 0, image); err != nil {
		return err
	}
	if err := oled.OLED.Draw(); err != nil {
		return err
	}
	return nil
}

// Clear ...
func (oled *OledDisplay) Clear() error {
	if err := oled.OLED.Clear(); err != nil {
		return err
	}
	return nil
}

// Close ...
func (oled *OledDisplay) Close() {
	oled.OLED.Clear()
	oled.OLED.Close()
}

// Off ...
func (oled *OledDisplay) Off() {
	oled.OLED.Clear()
	oled.OLED.Off()
}

func (oled *OledDisplay) drawText(text string, size float64, x, y int) (image.Image, error) {
	dst := image.NewRGBA(image.Rect(0, 0, oled.width, oled.height))
	draw.Draw(dst, dst.Bounds(), image.Transparent, image.Point{}, draw.Src)

	c := freetype.NewContext()
	c.SetDst(dst)
	c.SetClip(dst.Bounds())
	c.SetSrc(image.White)
	c.SetFont(oled.font)
	c.SetFontSize(size)

	if _, err := c.DrawString(text, freetype.Pt(x, y)); err != nil {
		return nil, err
	}

	return dst, nil
}
