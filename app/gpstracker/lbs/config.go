package lbs

import (
	"encoding/json"
	"io/ioutil"

	"github.com/shanghuiyang/rpi-devices/iot"
	"github.com/shanghuiyang/rpi-devices/util/geo"
)

// Config ...
type Config struct {
	ZoomInButtonPin  uint8          `json:"zoomInButtonPin"`
	ZoomOutButtonPin uint8          `json:"zoomOutButtonPin"`
	GPS              *GPSConfig     `json:"gps"`
	Display          *DisplayConfig `json:"display"`
	IOT              *IOTConfig     `json:"iot"`
	Tile             *TileConfig    `json:"tile"`
}

type GPSConfig struct {
	Dev        string              `json:"dev"`
	Baud       int                 `json:"baud"`
	DefaultLoc *geo.Point          `json:"defaultLoc"`
	Simulator  *GPSSimulatorConfig `json:"simulator"`
}

type GPSSimulatorConfig struct {
	Enable bool   `json:"enable"`
	Source string `json:"source"`
}

type DisplayConfig struct {
	Res    uint8 `json:"res"`
	Dc     uint8 `json:"dc"`
	Blk    uint8 `json:"blk"`
	Width  int   `json:"width"`
	Height int   `json:"height"`
}

type IOTConfig struct {
	Enable bool        `json:"enable"`
	Onenet *iot.Config `json:"onenet"`
}

type TileConfig struct {
	MinZoom             int      `json:"ninZoom"`
	MaxZoom             int      `json:"maxZoom"`
	DefaultZoom         int      `json:"defaultZoom"`
	TileProviders       []string `json:"tileProviders"`
	DefaultTileProvider string   `json:"defaultTileProvider"`
	CacheDir            string   `json:"cacheDir"`
	Online              bool     `json:"online"`
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
