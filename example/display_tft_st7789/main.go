package main

import (
	"bytes"
	"image/jpeg"
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
	data, err := ioutil.ReadFile("cat.jpg")
	if err != nil {
		log.Printf("open map.png error: %v", err)
		return
	}
	buf := bytes.NewBuffer(data)
	img, err := jpeg.Decode(buf)
	if err != nil {
		log.Printf("decode image error: %v", err)
		return
	}

	display, err := dev.NewST7789Display(res, dc, blk, width, height)
	if err != nil {
		log.Printf("failed to create an st display, error: %v", err)
		return
	}

	if err := display.Image(img); err != nil {
		log.Fatal(err)
	}
}
