package iot

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"

	"github.com/shanghuiyang/rpi-devices/base"
)

const (
	logTagOneNet = "onenet"
)

// OneNetCloud is the implement of Cloud
type OneNetCloud struct {
	token string
	api   string
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

// NewOneNetCloud ...
func NewOneNetCloud(cfg *base.OneNetConfig) *OneNetCloud {
	return &OneNetCloud{
		token: cfg.Token,
		api:   cfg.API,
	}
}

// Push ...
func (o *OneNetCloud) Push(v *Value) error {
	datapoint := OneNetData{
		Datastreams: []*Datastream{
			{
				ID: v.DeviceName,
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
