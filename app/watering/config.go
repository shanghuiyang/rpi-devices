package main

import (
	"encoding/json"
	"io/ioutil"

	"github.com/shanghuiyang/rpi-devices/iot"
)

type config struct {
	Pumps []*pumpConfig `json:"pumps"`
	Iot   *iotConfig    `json:"iot"`
}

type pumpConfig struct {
	Name        string `json:"name"`
	Pin         uint8  `json:"pin"`
	WateringAt  string `json:"wateringAt"`
	WateringSec int    `json:"wateringSec"`
	Button      uint8  `json:"button"`
}

type iotConfig struct {
	Enable bool        `json:"enable"`
	Onenet *iot.Config `json:"onenet"`
}

func loadConfig(file string) (*config, error) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err

	}
	var cfg config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
