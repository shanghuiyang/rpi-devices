package main

import (
	"encoding/json"
	"io/ioutil"

	"github.com/shanghuiyang/rpi-devices/iot"
)

type config struct {
	Stepper *stepperConfig `json:"stepper"`
	Button  uint8          `json:"button"`
	FeedAt  string         `json:"feedAt"`
	Iot     *iotConfig     `json:"iot"`
}

type stepperConfig struct {
	In1 uint8 `json:"in1"`
	In2 uint8 `json:"in2"`
	In3 uint8 `json:"in3"`
	In4 uint8 `json:"in4"`
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
