package main

import (
	"bytes"
	"log"
	"net/http"
	"time"

	"github.com/shanghuiyang/rpi-devices/dev"
)

const (
	onenetToken = "your_onenet_token"
	onenetAPI   = "http://api.heclouds.com/bindata?device_id=853487795&datastream_id=images"
)

func main() {
	cam := dev.NewMotionCamera()
	for {
		img, err := cam.Photo()
		if err != nil {
			log.Printf("failed to take phote from camera, error: %v", err)
			continue
		}
		go pushImage(img)
	}
}

func pushImage(img []byte) {
	buf := bytes.NewReader(img)
	req, err := http.NewRequest("POST", onenetAPI, buf)
	if err != nil {
		log.Printf("failed to new http request, error: %v", err)
		return
	}
	req.Header.Set("api-key", onenetToken)
	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("failed to send http request, error: %v", err)
		return
	}
	resp.Body.Close()
}
