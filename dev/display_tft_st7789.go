/*
ST7789Display is a driver for the tft lcd display module drived by ST7789 chip.
Please NOTE that you should disable SPI interface before using this driver!

Config Raspberry Pi:
1. $ sudo raspi-config
2. 	-> [5 Interface Options] -> [P4 SPI] -> [no] -> [ok]
3. $ sudo reboot now

Connect to Raspberry Pi:
 - VCC: any 3.3v or 5v pin
 - GND: any GND pin
 - SCL: GPIO 11(SPI0 SCLK)
 - SDA: GPIO 10(SPI0 MOSI)
 - RES: any data pin
 - DC:  any data pin
 - BLK: any data pin
*/

package dev

import (
	"errors"
	"image"
	"image/color"
	"image/draw"

	"github.com/stianeikeland/go-rpio/v4"
)

// ST7789Display is a driver for the tft lcd display module drived by ST7789 chip.
// It is an implement of Display interface.
type ST7789Display struct {
	res    rpio.Pin
	dc     rpio.Pin
	blk    rpio.Pin
	width  int
	height int
}

// NewST7789Display create a driver for the tft lcd display module drived by ST7789 chip.
// Note that you should disable SPI interface in raspi-config first!
func NewST7789Display(res, dc, blk uint8, width, height int) (*ST7789Display, error) {
	if err := rpio.SpiBegin(rpio.Spi0); err != nil {
		return nil, err
	}

	rpio.SpiMode(1, 0)
	rpio.SpiSpeed(40000000)

	display := &ST7789Display{}
	display.res = rpio.Pin(res)
	display.dc = rpio.Pin(dc)
	display.blk = rpio.Pin(blk)
	display.width = width
	display.height = height

	display.dc.Output()
	display.res.Output()
	display.blk.Output()
	display.blk.High()

	display.reset()
	display.init()

	return display, nil
}

// Image displays an image on the screen
func (display *ST7789Display) Image(img image.Image) error {
	display.setwindow()
	r := image.Rect(0, 0, display.width, display.height)
	dst := image.NewRGBA(r)
	draw.Draw(dst, r, img, r.Min, draw.Src)

	data := []byte{}
	for y := 0; y < dst.Bounds().Dy(); y++ {
		for x := 0; x < dst.Bounds().Dx(); x++ {
			c := dst.At(x, y).(color.RGBA)
			c565 := display.rgbaTo565(c)
			data = append(data, byte(c565>>8), byte(c565&0xFF))
		}
	}
	display.data(data...)
	return nil
}

// Text displays the text on the screen.
// NOTE: It isn't implemented. It is here just for implementing the Display interface.
// Please draw your text to an image first, and then use DisplayImage()
func (display *ST7789Display) Text(text string, x, y int) error {
	return errors.New("not implement")
}

// On turns the blacklight on
func (display *ST7789Display) On() error {
	display.blk.High()
	return nil
}

// On turns the blacklight off
func (display *ST7789Display) Off() error {
	display.blk.Low()
	return nil
}

// Clear clears the image
func (display *ST7789Display) Clear() error {
	// create a transparent image and display it
	r := image.Rect(0, 0, display.width, display.height)
	img := image.NewRGBA(r)
	draw.Draw(img, r, image.Transparent, r.Min, draw.Src)
	return display.Image(img)
}

// Close closes the module
func (display *ST7789Display) Close() error {
	_ = display.Clear()
	rpio.SpiEnd(rpio.Spi0)
	return nil
}

func (display *ST7789Display) reset() {
	display.res.High()
	delayMs(50)

	display.res.Low()
	delayMs(50)

	display.res.High()
	delayMs(50)
}

func (display *ST7789Display) init() {
	delayMs(10)
	display.command(0x11)
	delayMs(50)

	display.command(0x36)
	display.data(0x00)

	display.command(0x3A)
	display.data(0x05)

	display.command(0xB2)
	display.data(0x0C)
	display.data(0x0C)

	display.command(0xB7)
	display.data(0x35)

	display.command(0xBB)
	display.data(0x1A)

	display.command(0xC0)
	display.data(0x2C)

	display.command(0xC2)
	display.data(0x01)

	display.command(0xC3)
	display.data(0x0B)

	display.command(0xC4)
	display.data(0x20)

	display.command(0xC6)
	display.data(0x0F)

	display.command(0xD0)
	display.data(0xA4)
	display.data(0xA1)

	display.command(0x21)

	display.command(0xE0)
	display.data(0x00)
	display.data(0x19)
	display.data(0x1E)
	display.data(0x0A)
	display.data(0x09)
	display.data(0x15)
	display.data(0x3D)
	display.data(0x44)
	display.data(0x51)
	display.data(0x12)
	display.data(0x03)
	display.data(0x00)
	display.data(0x3F)
	display.data(0x3F)

	display.command(0xE1)
	display.data(0x00)
	display.data(0x18)
	display.data(0x1E)
	display.data(0x0A)
	display.data(0x09)
	display.data(0x25)
	display.data(0x3F)
	display.data(0x43)
	display.data(0x52)
	display.data(0x33)
	display.data(0x03)
	display.data(0x00)
	display.data(0x3F)
	display.data(0x3F)
	display.command(0x29)

	delayMs(50)
}

func (display *ST7789Display) setwindow() {
	minx, maxx := 0, display.width-1
	miny, maxy := 0, display.height-1
	display.command(0x2A)
	display.data(byte(minx >> 8))
	display.data(byte(minx)) // XSTART
	display.data(byte(maxx >> 8))
	display.data(byte(maxx)) // XEND
	display.command(0x2B)    // Row addr set
	display.data(byte(miny >> 8))
	display.data(byte(miny)) // YSTART
	display.data(byte(maxy >> 8))
	display.data(byte(maxy)) // YEND
	display.command(0x2C)    // write to RAM
}

func (display *ST7789Display) command(data ...byte) {
	display.dc.Low()
	rpio.SpiTransmit(data...)
}

func (display *ST7789Display) data(data ...byte) {
	display.dc.High()
	rpio.SpiTransmit(data...)
}

func (display *ST7789Display) rgbaTo565(c color.RGBA) uint16 {
	r, g, b, _ := c.RGBA()
	return uint16((r & 0xF800) + ((g & 0xFC00) >> 5) + ((b & 0xF800) >> 11))
}
