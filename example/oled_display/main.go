package main

import (
	"image"
	"image/draw"
	"io/ioutil"
	"log"
	"time"

	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"github.com/shanghuiyang/rpi-devices/dev"
	"github.com/shanghuiyang/rpi-devices/util"
)

const (
	width    = 128
	height   = 32
	fontFile = "casio-fx-9860gii.ttf"
)

func main() {
	oled, err := dev.NewOledDisplay(width, height)
	if err != nil {
		log.Printf("failed to create an oled, error: %v", err)
		return
	}

	util.WaitQuit(func() { oled.Close() })
	for {
		t := time.Now().Format("15:04:05")
		img, err := drawImage(t, 19, 0, 25)
		if err != nil {
			log.Printf("failed to draw image, error: %v", err)
			break
		}
		if err := oled.DisplayImage(img); err != nil {
			log.Printf("failed to display time, error: %v", err)
			break
		}
		time.Sleep(1 * time.Second)
	}
}

func drawImage(text string, size float64, x, y int) (image.Image, error) {
	data, err := ioutil.ReadFile(fontFile)
	if err != nil {
		return nil, err
	}
	font, err := truetype.Parse(data)
	if err != nil {
		return nil, err
	}

	r := image.Rect(0, 0, width, height)
	img := image.NewRGBA(r)
	draw.Draw(img, r, image.Transparent, r.Min, draw.Src)

	c := freetype.NewContext()
	c.SetDst(img)
	c.SetClip(r)
	c.SetSrc(image.White)
	c.SetFont(font)
	c.SetFontSize(size)

	if _, err := c.DrawString(text, freetype.Pt(x, y)); err != nil {
		return nil, err
	}

	return img, nil
}
