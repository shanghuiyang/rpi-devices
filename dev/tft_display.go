package dev

import (
	"image"
	"image/color"

	"github.com/stianeikeland/go-rpio/v4"
)

type TFTDisplay struct {
	res           rpio.Pin
	dc            rpio.Pin
	blk           rpio.Pin
	width, height int
}

func NewTFTDisplay(res, dc, blk uint8, width, height int) (*TFTDisplay, error) {
	if err := rpio.SpiBegin(rpio.Spi0); err != nil {
		return nil, err
	}

	rpio.SpiMode(1, 0)
	rpio.SpiSpeed(40000000)

	tft := &TFTDisplay{}
	tft.res = rpio.Pin(res)
	tft.dc = rpio.Pin(dc)
	tft.blk = rpio.Pin(blk)
	tft.width = width
	tft.height = height

	tft.dc.Output()
	tft.res.Output()
	tft.blk.Output()
	tft.blk.High()

	tft.reset()
	tft.init()

	return tft, nil
}

func (tft *TFTDisplay) reset() {
	tft.res.High()
	delayMs(50)

	tft.res.Low()
	delayMs(50)

	tft.res.High()
	delayMs(50)
}

func (tft *TFTDisplay) init() {
	delayMs(10)
	tft.command(0x11)
	delayMs(50)

	tft.command(0x36)
	tft.data(0x00)

	tft.command(0x3A)
	tft.data(0x05)

	tft.command(0xB2)
	tft.data(0x0C)
	tft.data(0x0C)

	tft.command(0xB7)
	tft.data(0x35)

	tft.command(0xBB)
	tft.data(0x1A)

	tft.command(0xC0)
	tft.data(0x2C)

	tft.command(0xC2)
	tft.data(0x01)

	tft.command(0xC3)
	tft.data(0x0B)

	tft.command(0xC4)
	tft.data(0x20)

	tft.command(0xC6)
	tft.data(0x0F)

	tft.command(0xD0)
	tft.data(0xA4)
	tft.data(0xA1)

	tft.command(0x21)

	tft.command(0xE0)
	tft.data(0x00)
	tft.data(0x19)
	tft.data(0x1E)
	tft.data(0x0A)
	tft.data(0x09)
	tft.data(0x15)
	tft.data(0x3D)
	tft.data(0x44)
	tft.data(0x51)
	tft.data(0x12)
	tft.data(0x03)
	tft.data(0x00)
	tft.data(0x3F)
	tft.data(0x3F)

	tft.command(0xE1)
	tft.data(0x00)
	tft.data(0x18)
	tft.data(0x1E)
	tft.data(0x0A)
	tft.data(0x09)
	tft.data(0x25)
	tft.data(0x3F)
	tft.data(0x43)
	tft.data(0x52)
	tft.data(0x33)
	tft.data(0x03)
	tft.data(0x00)
	tft.data(0x3F)
	tft.data(0x3F)
	tft.command(0x29)

	delayMs(50)
}

func (tft *TFTDisplay) Clear() {

}

func (tft *TFTDisplay) Close() {
	rpio.SpiEnd(rpio.Spi0)
}

func (tft *TFTDisplay) Display(img image.Image) error {
	tft.setwindow()
	data := []byte{}
	for y := 0; y < img.Bounds().Dy(); y++ {
		for x := 0; x < img.Bounds().Dx(); x++ {
			c := img.At(x, y)
			r, g, b, a := c.RGBA()
			c565 := tft.rgbaTo565(color.RGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: uint8(a)})
			data = append(data, byte(c565>>8))
			data = append(data, byte(c565&0xFF))
		}
	}
	tft.data(data...)
	return nil
}

func (tft *TFTDisplay) setwindow() {
	minx, maxx := 0, tft.width-1
	miny, maxy := 0, tft.height-1
	tft.command(0x2A)
	tft.data(byte(minx >> 8))
	tft.data(byte(minx)) // XSTART
	tft.data(byte(maxx >> 8))
	tft.data(byte(maxx)) // XEND
	tft.command(0x2B)    // Row addr set
	tft.data(byte(miny >> 8))
	tft.data(byte(miny)) // YSTART
	tft.data(byte(maxy >> 8))
	tft.data(byte(maxy)) // YEND
	tft.command(0x2C)    // write to RAM
}

func (tft *TFTDisplay) command(data ...byte) {
	tft.dc.Low()
	rpio.SpiTransmit(data...)
}

func (tft *TFTDisplay) data(data ...byte) {
	tft.dc.High()
	rpio.SpiTransmit(data...)
}

func (tft *TFTDisplay) rgbaTo565(c color.RGBA) uint16 {
	r, g, b, _ := c.RGBA()
	return uint16((r & 0xF800) + ((g & 0xFC00) >> 5) + ((b & 0xF800) >> 11))
}
