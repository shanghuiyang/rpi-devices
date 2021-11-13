package main

import (
	"bytes"
	"encoding/json"
	"errors"
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

func main() {
	streamer, err := util.NewStreamer(streamerURL)
	if err != nil {
		log.Fatalf("failed to create streamer, error: %v", err)
	}

	for {
		idx, err := getImageIndex()
		if err != nil {
			log.Printf("failed to get image index, error: %v", err)
			continue
		}
		log.Printf("%v", idx)

		data, err := cloud.GetFile(idx)
		if err != nil {
			log.Printf("failed to get image data, error: %v", err)
			continue
		}

		img, _, err := image.Decode(bytes.NewReader(data))
		if err != nil {
			log.Printf("decode image error: %v", err)
			continue
		}
		buf := new(bytes.Buffer)
		if err := jpeg.Encode(buf, img, nil); err != nil {
			log.Printf("encode image error: %v", err)
			continue
		}
		streamer.Push(buf.Bytes())
	}

}

func getImageIndex() (string, error) {
	params := map[string]interface{}{
		"datastream_id": "images",
		"limit":         1,
	}
	result, err := cloud.Get(params)
	if err != nil {
		return "", err
	}

	var data iot.OnenetData
	if err := json.Unmarshal(result, &data); err != nil {
		return "", err
	}
	if len(data.Datastreams) == 0 {
		return "", errors.New("empty data")
	}
	images, ok := data.Datastreams[0].Datapoints[0].Value.(map[string]interface{})
	if !ok {
		return "", errors.New("can't convert value to map[string]interface{}")
	}
	idx, ok := images["index"]
	if !ok {
		return "", err
	}

	s, ok := idx.(string)
	if !ok {
		return "", errors.New("can't convert interface to string")
	}

	return s, nil
}
