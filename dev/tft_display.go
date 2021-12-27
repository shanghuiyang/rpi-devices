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

func NewTFTDisplay(res, dc, blk uint8) (*TFTDisplay, error) {
	if err := rpio.SpiBegin(rpio.Spi0); err != nil {
		return nil, err
	}

	rpio.SpiMode(1, 0)
	// rpio.SpiSpeed(40000000)

	tft := &TFTDisplay{}
	tft.res = rpio.Pin(res)
	tft.dc = rpio.Pin(dc)
	tft.blk = rpio.Pin(blk)
	tft.width = 240
	tft.height = 240

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
	tft.sendCmd(0x11)
	delayMs(50)

	tft.sendCmd(0x36)
	tft.sendData(0x00)

	tft.sendCmd(0x3A)
	tft.sendData(0x05)

	tft.sendCmd(0xB2)
	tft.sendData(0x0C)
	tft.sendData(0x0C)

	tft.sendCmd(0xB7)
	tft.sendData(0x35)

	tft.sendCmd(0xBB)
	tft.sendData(0x1A)

	tft.sendCmd(0xC0)
	tft.sendData(0x2C)

	tft.sendCmd(0xC2)
	tft.sendData(0x01)

	tft.sendCmd(0xC3)
	tft.sendData(0x0B)

	tft.sendCmd(0xC4)
	tft.sendData(0x20)

	tft.sendCmd(0xC6)
	tft.sendData(0x0F)

	tft.sendCmd(0xD0)
	tft.sendData(0xA4)
	tft.sendData(0xA1)

	tft.sendCmd(0x21)

	tft.sendCmd(0xE0)
	tft.sendData(0x00)
	tft.sendData(0x19)
	tft.sendData(0x1E)
	tft.sendData(0x0A)
	tft.sendData(0x09)
	tft.sendData(0x15)
	tft.sendData(0x3D)
	tft.sendData(0x44)
	tft.sendData(0x51)
	tft.sendData(0x12)
	tft.sendData(0x03)
	tft.sendData(0x00)
	tft.sendData(0x3F)
	tft.sendData(0x3F)

	tft.sendCmd(0xE1)
	tft.sendData(0x00)
	tft.sendData(0x18)
	tft.sendData(0x1E)
	tft.sendData(0x0A)
	tft.sendData(0x09)
	tft.sendData(0x25)
	tft.sendData(0x3F)
	tft.sendData(0x43)
	tft.sendData(0x52)
	tft.sendData(0x33)
	tft.sendData(0x03)
	tft.sendData(0x00)
	tft.sendData(0x3F)
	tft.sendData(0x3F)
	tft.sendCmd(0x29)

	delayMs(50)
}

func (tft *TFTDisplay) Clear() {

}

func (tft *TFTDisplay) Close() {
	rpio.SpiEnd(rpio.Spi0)
}

func (tft *TFTDisplay) Display(img image.Image) error {
	minx, maxx := 0, tft.width-1
	miny, maxy := 0, tft.height-1
	tft.sendCmd(0x2A)
	tft.sendData(byte(minx >> 8))
	tft.sendData(byte(minx)) // XSTART
	tft.sendData(byte(maxx >> 8))
	tft.sendData(byte(maxx)) // XEND
	tft.sendCmd(0x2B)        // Row addr set
	tft.sendData(byte(miny >> 8))
	tft.sendData(byte(miny)) // YSTART
	tft.sendData(byte(maxy >> 8))
	tft.sendData(byte(maxy)) // YEND
	tft.sendCmd(0x2C)        // write to RAM

	data := []byte{}
	for y := 0; y < img.Bounds().Dy(); y++ {
		for x := 0; x < img.Bounds().Dx(); x++ {
			c := img.At(x, y)
			r, g, b, a := c.RGBA()
			c565 := rgbaTo565(color.RGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: uint8(a)})
			data = append(data, byte(c565>>8))
			data = append(data, byte(c565&0xFF))
		}
	}
	tft.sendData(data...)
	return nil
}

func (tft *TFTDisplay) sendCmd(data ...byte) {
	tft.dc.Low()
	rpio.SpiTransmit(data...)
}

func (tft *TFTDisplay) sendData(data ...byte) {
	tft.dc.High()
	rpio.SpiTransmit(data...)
}

func rgbaTo565(c color.RGBA) uint16 {
	r, g, b, _ := c.RGBA()
	return uint16((r & 0xF800) + ((g & 0xFC00) >> 5) + ((b & 0xF800) >> 11))
}
