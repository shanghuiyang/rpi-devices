package iot

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"
)

// Onenet is the implement of Cloud
type Onenet struct {
	token string
	api   string
}

type onenetData struct {
	OnenetStreams []*onenetStream `json:"datastreams"`
}

type onenetStream struct {
	ID         string       `json:"id"`
	Datapoints []*datapoint `json:"datapoints"`
}

type datapoint struct {
	Value interface{} `json:"value"`
}

// NewOnenet ...
func NewOnenet(cfg *Config) *Onenet {
	return &Onenet{
		token: cfg.Token,
		api:   cfg.API,
	}
}

// Push ...
func (o *Onenet) Push(v *Value) error {
	datapoint := onenetData{
		OnenetStreams: []*onenetStream{
			{
				ID: v.Device,
				Datapoints: []*datapoint{
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
	if err != nil {
		return err
	}
	req.Header.Set("api-key", o.token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		Timeout: 15 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}
