package main

import (
	"bytes"
	"encoding/json"
	"image"
	"image/jpeg"
	"log"

	"github.com/shanghuiyang/rpi-devices/iot"
	"github.com/shanghuiyang/rpi-devices/util"
)

const (
	streamerURL = ":8088/monitor"
	onenetToken = "your_onenet_token"
	onenetAPI   = "http://api.heclouds.com/devices/853487795/datapoints"
)

var cloud = iot.NewOnenet(&iot.Config{
	API:   onenetAPI,
	Token: onenetToken,
})

var (
	// !NOTE: the channel size should be 1 or will result to obvious delays.
	chImageIdx  = make(chan string, 1)
	chImageData = make(chan []byte, 1)
)

func main() {
	go getImageIndex()
	go getFile()
	go pushImg()
	select {}
}

func getImageIndex() {
	logTag := "getImageIndex"
	for {
		params := map[string]interface{}{
			"datastream_id": "images",
			"limit":         1,
		}
		result, err := cloud.Get(params)
		if err != nil {
			log.Printf("[%v]Get() failed, error: %v", logTag, err)
			continue
		}

		var data iot.OnenetData
		if err := json.Unmarshal(result, &data); err != nil {
			continue
		}
		if len(data.Datastreams) == 0 {
			log.Printf("[%v]empty", logTag)
			continue
		}
		images, ok := data.Datastreams[0].Datapoints[0].Value.(map[string]interface{})
		if !ok {
			log.Printf("[%v]can't convert value to map[string]interface{}", logTag)
			continue
		}
		idx, ok := images["index"]
		if !ok {
			log.Printf("[%v]invalid result", logTag)
			continue
		}

		s, ok := idx.(string)
		if !ok {
			log.Printf("[%v]can't convert interface to string", logTag)
			continue
		}

		chImageIdx <- s
	}
}

func getFile() {
	logTag := "getFile"
	for idx := range chImageIdx {
		img, err := cloud.GetFile(idx)
		if err != nil {
			log.Printf("[%v]failed to get image data, error: %v", logTag, err)
			continue
		}
		chImageData <- img
	}
}

func pushImg() {
	logTag := "pushImg"
	streamer, err := util.NewStreamer(streamerURL)
	if err != nil {
		log.Fatalf("[%v]failed to create streamer, error: %v", logTag, err)
	}

	for data := range chImageData {
		img, _, err := image.Decode(bytes.NewReader(data))
		if err != nil {
			log.Printf("[%v]decode image error: %v", logTag, err)
			continue
		}
		buf := new(bytes.Buffer)
		if err := jpeg.Encode(buf, img, nil); err != nil {
			log.Printf("[%v]encode image error: %v", logTag, err)
			continue
		}
		streamer.Push(buf.Bytes())
	}
}
