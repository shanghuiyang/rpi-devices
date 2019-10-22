package iot

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/shanghuiyang/rpi-devices/base"
)

const (
	logTagWsn = "wsn"
)

// WsnCloud is the implement of Cloud
type WsnCloud struct {
	token string
	api   string
}

// NewWsnClound ...
func NewWsnClound(cfg *base.WsnConfig) *WsnCloud {
	return &WsnCloud{
		token: cfg.Token,
		api:   cfg.API,
	}
}

// Push ...
func (w *WsnCloud) Push(v *Value) error {
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
			"id":    {v.DeviceName},
			"lat":   {fmt.Sprintf("%v", pt.Lat)},
			"lng":   {fmt.Sprintf("%v", pt.Lon)},
			"speed": {"30"},
		}
	} else {
		formData = url.Values{
			"ak":    {w.token},
			"id":    {v.DeviceName},
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
