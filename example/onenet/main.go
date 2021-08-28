package main

import (
	"encoding/json"
	"log"

	"github.com/shanghuiyang/rpi-devices/dev"
	"github.com/shanghuiyang/rpi-devices/iot"
	"github.com/shanghuiyang/rpi-devices/util"
)

const (
	ledPin = 12

	onenetToken = "your_onenet_token"
	onenetAPI   = "http://api.heclouds.com/devices/540381180/datapoints"
)

func main() {
	led := dev.NewLedImp(ledPin)
	o := iot.NewOnenet(&iot.Config{
		Token: onenetToken,
		API:   onenetAPI,
	})

	for {
		util.DelayMs(1000)
		params := map[string]interface{}{
			"datastream_id": "light",
			"limit":         1,
		}
		result, err := o.Get(params)
		if err != nil {
			log.Printf("failed to get data from onenet, error: %v", err)
			continue
		}
		var data iot.OnenetData
		if err := json.Unmarshal(result, &data); err != nil {
			log.Printf("failed to unmarshal data, err: %v", err)
			continue
		}
		if len(data.Datastreams) == 0 {
			log.Printf("empty data")
			continue
		}
		turnon := data.Datastreams[0].Datapoints[0].Value.(float64) == 1
		if turnon {
			led.On()
		} else {
			led.Off()
		}
	}
}
