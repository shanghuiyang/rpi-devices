package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/shanghuiyang/rpi-devices/iot"
	"github.com/shanghuiyang/rpi-devices/util"
)

const (
	streamerURL             = ":8088/monitor"
	onenetToken             = "your_onenet_token"
	onenetAPI               = "http://api.heclouds.com/devices/853487795/datapoints"
	onenetGetFileAPIPattern = "your_onenet_get_file_api_pattern_%v"
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
		imgidx, err := getImageIndex()
		if err != nil {
			log.Printf("failed to get image index, error: %v", err)
			continue
		}
		log.Printf("%v", imgidx)

		api := fmt.Sprintf(onenetGetFileAPIPattern, imgidx)
		imgdata, err := getFile(api)
		if err != nil {
			log.Printf("failed to get image data, error: %v", err)
			continue
		}

		img, _, err := image.Decode(bytes.NewReader(imgdata))
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
	imgidx, ok := images["index"]
	if !ok {
		return "", err
	}

	imgidxStr, ok := imgidx.(string)
	if !ok {
		return "", errors.New("can't convert interface to string")
	}

	return imgidxStr, nil
}

func getFile(api string) ([]byte, error) {
	req, err := http.NewRequest("GET", api, nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{
		Timeout: 3 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return data, nil
}
