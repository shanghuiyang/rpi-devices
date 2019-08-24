package iotclouds

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/shanghuiyang/rpi-devices/base"
)

const (
	logTagWsn = "wsn"
)

// WsnCloud is the implement of IOTCloud
type WsnCloud struct {
	token   string
	api     string
	devices map[string]string // key: device name; value: device id
}

// NewWsnClound ...
func NewWsnClound(cfg *base.WsnConfig) *WsnCloud {
	return &WsnCloud{
		token:   cfg.Token,
		api:     cfg.API,
		devices: cfg.Devices,
	}
}

// Start ...
func (w *WsnCloud) Start() {
	log.Printf("[%v]start working", logTagWsn)
	for v := range chIoTCloud {
		if err := w.upload(v); err != nil {
			log.Printf("[%v]faied to push data to iot cloud, error: %v", logTagWsn, err)
		}
	}
	return
}

// Push ...
func (w *WsnCloud) Push(v *IoTValue) {
	chIoTCloud <- v
}

// Push ...
func (w *WsnCloud) upload(v *IoTValue) error {
	deviceID, ok := w.devices[v.DeviceName]
	if !ok {
		return fmt.Errorf(`device "%v" not found`, v.DeviceName)
	}

	var formData url.Values
	api := w.api
	if v.DeviceName == "gps" {
		api = strings.Replace(w.api, "numerical", "gps", -1)
		pt, ok := v.Value.(*base.Point)
		if !ok {
			return fmt.Errorf("failed to convert value to point")
		}
		formData = url.Values{
			"ak":    {w.token},
			"id":    {deviceID},
			"lat":   {fmt.Sprintf("%v", pt.Lat)},
			"lng":   {fmt.Sprintf("%v", pt.Lon)},
			"speed": {"30"},
		}
	} else {
		formData = url.Values{
			"ak":    {w.token},
			"id":    {deviceID},
			"value": {fmt.Sprintf("%v", v.Value)},
		}
	}

	resp, err := http.PostForm(api, formData)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}
