/*
Package dev ...

OLED is the driver of an oled screen.
Please NOTE that current version only supports the oled module with ssd1306 driver.

connect to raspberry pi:
VCC: pin 1 or any 3.3v pin
GND: pin 9 or and GND pin
SDA: pin 3 (SDA)
SCL: pin 5 (SCL)

		+--------------------+
		|        OLED        |
		|                    |
		+--+----+----+----+--+
		   |    |    |    |
		  GND  VCC  SCL  SDA
		   |    |    |    |
		   |    |    |    |
		   |    |    |    |            +-----------+
		   |    +----|----|------------+ * 1   2 o |
		   |         |    +------------| * 3     o |
		   |         +-----------------+ * 5     o |
		   |                           | o       o |
		   +---------------------------+ * 9     o |
		                               | o       o |
		                               | o       o |
		                               | o       o |
		                               | o       o |
		                               | o       o |
		                               | o       o |
		                               | o       o |
		                               | o       o |
		                               | o       o |
		                               | o       o |
		                               | o       o |
		                               | o       o |
		                               | o       o |
		                               | o       o |
		                               | o 39 40 o |
									   +-----------+


Config Your Pi:
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

// OLED ...
type OLED struct {
	oled   *monochromeoled.OLED
	width  int
	height int
	font   *truetype.Font
}

// NewOLED ...
func NewOLED(width, heigth int) (*OLED, error) {
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
	return &OLED{
		oled:   oled,
		width:  width,
		height: heigth,
		font:   font,
	}, nil
}

// Display ...
func (o *OLED) Display(text string, fontSize float64, x, y int) error {
	image, err := o.drawText(text, fontSize, x, y)
	if err != nil {
		return err
	}
	if err := o.oled.SetImage(0, 0, image); err != nil {
		return err
	}
	if err := o.oled.Draw(); err != nil {
		return err
	}
	return nil
}

// Clear ...
func (o *OLED) Clear() error {
	if err := o.oled.Clear(); err != nil {
		return err
	}
	return nil
}

// Close ...
func (o *OLED) Close() {
	o.oled.Clear()
	o.oled.Close()
}

// Off ...
func (o *OLED) Off() {
	o.oled.Clear()
	o.oled.Off()
}

func (o *OLED) drawText(text string, size float64, x, y int) (image.Image, error) {
	dst := image.NewRGBA(image.Rect(0, 0, o.width, o.height))
	draw.Draw(dst, dst.Bounds(), image.Transparent, image.ZP, draw.Src)

	c := freetype.NewContext()
	c.SetDst(dst)
	c.SetClip(dst.Bounds())
	c.SetSrc(image.White)
	c.SetFont(o.font)
	c.SetFontSize(size)

	if _, err := c.DrawString(text, freetype.Pt(x, y)); err != nil {
		return nil, err
	}

	return dst, nil
}
