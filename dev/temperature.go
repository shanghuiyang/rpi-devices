package dev

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

var (
	tempFile = "/sys/bus/w1/devices/28-d8baf71d64ff/w1_slave"
)

// Temperature ...
type Temperature struct {
}

// NewTemperature ...
func NewTemperature() *Temperature {
	return &Temperature{}
}

// GetTemperature ...
// temperature file:
// ------------------------------------------
// ca 01 55 00 7f ff 0c 10 bf : crc=bf YES
// ca 01 55 00 7f ff 0c 10 bf t=28625
// ------------------------------------------
func (t *Temperature) GetTemperature() (float32, error) {
	data, err := ioutil.ReadFile(tempFile)
	if err != nil {
		return 0, err
	}
	raw := string(data)

	idx := strings.LastIndex(raw, "t=")
	if idx < 0 {
		return 0, fmt.Errorf("can't find 't='")
	}
	c, err := strconv.ParseFloat(raw[idx+2:idx+7], 32)
	if err != nil {
		return 0, fmt.Errorf("bad data")
	}
	return float32(c / 1000), nil
}
