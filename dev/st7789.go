package dev

import (
	"image"
	"image/color"
	"image/draw"

	"github.com/stianeikeland/go-rpio/v4"
)

type ST7789 struct {
	res    rpio.Pin
	dc     rpio.Pin
	blk    rpio.Pin
	width  int
	height int
}

// NewST7789 ...
func NewST7789(res, dc, blk uint8, width, height int) (*ST7789, error) {
	if err := rpio.SpiBegin(rpio.Spi0); err != nil {
		return nil, err
	}

	rpio.SpiMode(1, 0)
	rpio.SpiSpeed(40000000)

	st := &ST7789{}
	st.res = rpio.Pin(res)
	st.dc = rpio.Pin(dc)
	st.blk = rpio.Pin(blk)
	st.width = width
	st.height = height

	st.dc.Output()
	st.res.Output()
	st.blk.Output()
	st.blk.High()

	st.reset()
	st.init()

	return st, nil
}

// Display displays the img image on the module
func (st *ST7789) Display(img image.Image) error {
	st.setwindow()
	r := image.Rect(0, 0, st.width, st.height)
	dst := image.NewRGBA(r)
	draw.Draw(dst, r, img, r.Min, draw.Src)

	data := []byte{}
	for y := 0; y < dst.Bounds().Dy(); y++ {
		for x := 0; x < dst.Bounds().Dx(); x++ {
			c := dst.At(x, y).(color.RGBA)
			c565 := st.rgbaTo565(c)
			data = append(data, byte(c565>>8), byte(c565&0xFF))
		}
	}
	st.data(data...)
	return nil
}

// Clear clears the image
func (st *ST7789) Clear() {
	r := image.Rect(0, 0, st.width, st.height)
	img := image.NewRGBA(r)
	draw.Draw(img, r, image.Transparent, r.Min, draw.Src)
	st.Display(img)
}

// On turns the blacklight on
func (st *ST7789) On() {
	st.blk.High()
}

// On turns the blacklight off
func (st *ST7789) Off() {
	st.blk.Low()
}

// Close closes the module
func (st *ST7789) Close() {
	rpio.SpiEnd(rpio.Spi0)
}

func (st *ST7789) reset() {
	st.res.High()
	delayMs(50)

	st.res.Low()
	delayMs(50)

	st.res.High()
	delayMs(50)
}

func (st *ST7789) init() {
	delayMs(10)
	st.command(0x11)
	delayMs(50)

	st.command(0x36)
	st.data(0x00)

	st.command(0x3A)
	st.data(0x05)

	st.command(0xB2)
	st.data(0x0C)
	st.data(0x0C)

	st.command(0xB7)
	st.data(0x35)

	st.command(0xBB)
	st.data(0x1A)

	st.command(0xC0)
	st.data(0x2C)

	st.command(0xC2)
	st.data(0x01)

	st.command(0xC3)
	st.data(0x0B)

	st.command(0xC4)
	st.data(0x20)

	st.command(0xC6)
	st.data(0x0F)

	st.command(0xD0)
	st.data(0xA4)
	st.data(0xA1)

	st.command(0x21)

	st.command(0xE0)
	st.data(0x00)
	st.data(0x19)
	st.data(0x1E)
	st.data(0x0A)
	st.data(0x09)
	st.data(0x15)
	st.data(0x3D)
	st.data(0x44)
	st.data(0x51)
	st.data(0x12)
	st.data(0x03)
	st.data(0x00)
	st.data(0x3F)
	st.data(0x3F)

	st.command(0xE1)
	st.data(0x00)
	st.data(0x18)
	st.data(0x1E)
	st.data(0x0A)
	st.data(0x09)
	st.data(0x25)
	st.data(0x3F)
	st.data(0x43)
	st.data(0x52)
	st.data(0x33)
	st.data(0x03)
	st.data(0x00)
	st.data(0x3F)
	st.data(0x3F)
	st.command(0x29)

	delayMs(50)
}

func (st *ST7789) setwindow() {
	minx, maxx := 0, st.width-1
	miny, maxy := 0, st.height-1
	st.command(0x2A)
	st.data(byte(minx >> 8))
	st.data(byte(minx)) // XSTART
	st.data(byte(maxx >> 8))
	st.data(byte(maxx)) // XEND
	st.command(0x2B)    // Row addr set
	st.data(byte(miny >> 8))
	st.data(byte(miny)) // YSTART
	st.data(byte(maxy >> 8))
	st.data(byte(maxy)) // YEND
	st.command(0x2C)    // write to RAM
}

func (st *ST7789) command(data ...byte) {
	st.dc.Low()
	rpio.SpiTransmit(data...)
}

func (st *ST7789) data(data ...byte) {
	st.dc.High()
	rpio.SpiTransmit(data...)
}

func (st *ST7789) rgbaTo565(c color.RGBA) uint16 {
	r, g, b, _ := c.RGBA()
	return uint16((r & 0xF800) + ((g & 0xFC00) >> 5) + ((b & 0xF800) >> 11))
}
