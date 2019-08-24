package iotclouds

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/shanghuiyang/rpi-devices/base"
)

const (
	logTagOneNet = "onenet"
)

// OneNetCloud is the implement of IOTCloud
type OneNetCloud struct {
	token   string
	api     string
	devices map[string]string // key: device name; value: device id
}

// OneNetData ...
type OneNetData struct {
	Datastreams []*Datastream `json:"datastreams"`
}

// Datastream ...
type Datastream struct {
	ID         string       `json:"id"`
	Datapoints []*Datapoint `json:"datapoints"`
}

// Datapoint ...
type Datapoint struct {
	Value interface{} `json:"value"`
}

// NewOneNetClound ...
func NewOneNetClound(cfg *base.OneNetConfig) *OneNetCloud {
	return &OneNetCloud{
		token:   cfg.Token,
		api:     cfg.API,
		devices: cfg.Devices,
	}
}

// Start ...
func (o *OneNetCloud) Start() {
	log.Printf("[%v]start working", logTagOneNet)
	for v := range chIoTCloud {
		go func(v *IoTValue) {
			if err := o.upload(v); err != nil {
				log.Printf("[%v]faied to push data to iot cloud, error: %v", logTagOneNet, err)
			}
		}(v)
	}
	return
}

// Push ...
func (o *OneNetCloud) Push(v *IoTValue) {
	chIoTCloud <- v
}

// Upload ...
func (o *OneNetCloud) upload(v *IoTValue) error {
	deviceID, ok := o.devices[v.DeviceName]
	if !ok {
		return fmt.Errorf(`device "%v" not found`, v.DeviceName)
	}
	datapoint := OneNetData{
		Datastreams: []*Datastream{
			{
				ID: deviceID,
				Datapoints: []*Datapoint{
					{
						Value: v.Value,
					},
				},
			},
		},
	}
	data, err := json.Marshal(datapoint)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", o.api, bytes.NewBuffer(data))
	req.Header.Set("api-key", o.token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}
