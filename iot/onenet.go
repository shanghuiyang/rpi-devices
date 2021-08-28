package iot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

// Onenet is the implement of Cloud
type Onenet struct {
	token string
	api   string
}

type onenetData struct {
	Datastreams []*datastream `json:"datastreams"`
}

type datastream struct {
	ID         string       `json:"id"`
	Datapoints []*datapoint `json:"datapoints"`
}

type datapoint struct {
	At    string      `json:"at"`
	Value interface{} `json:"value"`
}

type OnenetDataResponse struct {
	Errno int64       `json:"errno"`
	Error string      `json:"error"`
	Data  *OnenetData `json:"data"`
}

type OnenetData struct {
	Count       int64         `json:"count"`
	Cursor      string        `json:"cursor"`
	Datastreams []*datastream `json:"datastreams"`
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
		Datastreams: []*datastream{
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

func (o *Onenet) Get(params map[string]interface{}) ([]byte, error) {
	api := o.api
	if len(params) > 0 {
		api += "?"
	}

	for k, v := range params {
		api += fmt.Sprintf("%v=%v&", k, v)
	}
	api = strings.TrimRight(api, "&")

	req, err := http.NewRequest("GET", api, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("api-key", o.token)

	client := &http.Client{
		Timeout: 15 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read resp body, err: %v", err)
	}
	var result OnenetDataResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal resp, err: %v", err)
	}
	if result.Errno > 0 {
		return nil, fmt.Errorf("error no: %v, message: %v", result.Errno, result.Error)
	}
	return json.Marshal(result.Data)
}
