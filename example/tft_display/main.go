package main

import (
	"bytes"
	"image/png"
	"io/ioutil"
	"log"

	"github.com/shanghuiyang/rpi-devices/dev"
)

const (
	res = 22
	dc  = 17
	blk = 27

	width  = 240
	height = 240
)

func main() {
	data, err := ioutil.ReadFile("map.png")
	if err != nil {
		log.Printf("open map.png error: %v", err)
		return
	}
	buf := bytes.NewBuffer(data)
	img, err := png.Decode(buf)
	if err != nil {
		log.Printf("decode image error: %v", err)
		return
	}

	tft, err := dev.NewTFTDisplay(res, dc, blk, width, height)
	if err != nil {
		log.Printf("failed to create an tft display, error: %v", err)
		return
	}
	defer tft.Close()

	if err := tft.Display(img); err != nil {
		log.Fatal(err)
	}
}
