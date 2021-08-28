package main

import (
	"encoding/json"
	"log"
	"time"

	"github.com/shanghuiyang/rpi-devices/dev"
	"github.com/shanghuiyang/rpi-devices/iot"
	"github.com/shanghuiyang/rpi-devices/util"
	"github.com/stianeikeland/go-rpio/v4"
)

const (
	d0 = 16
	d1 = 20
	d2 = 21
	d3 = 6

	localLedPin = 26
	cloudLedPin = 12

	butonAchannel = 3
	butonBchannel = 2
	butonCchannel = 1
	butonDchannel = 0

	onenetToken = "your_onenet_token"
	onenetAPI   = "http://api.heclouds.com/devices/540381180/datapoints"

	// use this rpio as 3.3v pin
	// if all 3.3v pins were used
	pin33v = 5
)

var light *rlight

type rlight struct {
	localLed dev.Led
	cloudLed dev.Led
	cloud    iot.Cloud
	rf       dev.RFReceiver
	on       bool
}

func main() {
	p33v := rpio.Pin(pin33v)
	p33v.Output()
	p33v.High()

	localLed := dev.NewLedImp(localLedPin)
	cloudLed := dev.NewLedImp(cloudLedPin)
	r := dev.NewRX480E4(d0, d1, d2, d3)
	cloud := iot.NewOnenet(&iot.Config{
		API:   onenetAPI,
		Token: onenetToken,
	})
	light = &rlight{
		localLed: localLed,
		cloudLed: cloudLed,
		cloud:    cloud,
		rf:       r,
		on:       false,
	}

	util.WaitQuit(func() {
		localLed.Off()
		cloudLed.Off()
	})

	go light.toggledByRF()
	go light.toggledByCloud()

	select {}

}

func (r *rlight) toggledByRF() {
	for {
		if r.rf.Received(butonAchannel) {
			go light.toggle()
			continue
		}
		if r.rf.Received(butonBchannel) {
			go light.toggle()
			continue
		}
		if r.rf.Received(butonCchannel) {
			go light.toggle()
			continue
		}
		if r.rf.Received(butonDchannel) {
			go light.toggle()
			continue
		}
		time.Sleep(300 * time.Millisecond)
	}
}

func (r *rlight) toggledByCloud() {
	for {
		util.DelayMs(1000)
		params := map[string]interface{}{
			"datastream_id": "light",
			"limit":         1,
		}
		result, err := r.cloud.Get(params)
		if err != nil {
			log.Printf("failed to get data from onenet, error: %v", err)
			continue
		}
		switch r.cloud.(type) {
		case *iot.Onenet:
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
				r.cloudLed.On()
			} else {
				r.cloudLed.Off()
			}
		default:
			log.Printf("not implement the cloud but onenet")
			continue
		}

	}
}

func (r *rlight) toggle() {
	if r.on {
		r.localLed.Off()
		r.on = false
		log.Printf("[rlight]light off")
	} else {
		r.localLed.On()
		r.on = true
		log.Printf("[rlight]light on")
	}
}
