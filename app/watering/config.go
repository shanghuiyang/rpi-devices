package main

import (
	"encoding/json"
	"io/ioutil"

	"github.com/shanghuiyang/rpi-devices/iot"
)

type config struct {
	Gardeners []*gardenercfg `json:"gardeners"`
	Button    uint8          `json:"button"`
	Iot       *iotcfg        `json:"iot"`
}

type gardenercfg struct {
	Name       string `json:"name"`
	WorkAt     string `json:"workAt"`
	WorkingSec int    `json:"workingSec"`
	Enabled    bool   `json:"enabled"`
	Relay      uint8  `json:"relay"`
}

type iotcfg struct {
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
