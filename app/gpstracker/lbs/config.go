package lbs

import (
	"encoding/json"
	"io/ioutil"

	"github.com/shanghuiyang/rpi-devices/iot"
	"github.com/shanghuiyang/rpi-devices/util/geo"
)

// Config ...
type Config struct {
	ZoomInButtonPin  uint8       `json:"zoomInButtonPin"`
	ZoomOutButtonPin uint8       `json:"zoomOutButtonPin"`
	GPS              *GPSConfig  `json:"gps"`
	IOT              *iot.Config `json:"iot"`
	Tile             *TileConfig `json:"tile"`
	Online           bool        `json:"online"`
	DefaultLocation  *geo.Point  `json:"defaultLocation"`
}

type GPSConfig struct {
	Dev  string `json:"dev"`
	Baud int    `json:"baud"`
}

type BaiduKey struct {
	APIKey    string `json:"apiKey"`
	SecretKey string `json:"secretKey"`
}

type TileConfig struct {
	MinZoom             int      `json:"ninZoom"`
	MaxZoom             int      `json:"maxZoom"`
	DefaultZoom         int      `json:"defaultZoom"`
	TileProviders       []string `json:"tileProviders"`
	DefaultTileProvider string   `json:"defaultTileProvider"`
	CacheDir            string   `json:"cacheDir"`
}

func LoadConfig(file string) (*Config, error) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err

	}
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
