package server

import (
	"encoding/json"
	"io/ioutil"

	"github.com/shanghuiyang/rpi-devices/util/geo"
)

// Config ...
type Config struct {
	LedPin         uint8                `json:"led"`
	SG90DataPin    uint8                `json:"sg90"`
	BuzzerPin      uint8                `json:"buzzer"`
	L298N          *L298NConfig         `json:"l298n"`
	US100          *US100Config         `json:"us100"`
	GY25           *GY25Config          `json:"gy25"`
	Joystick       *JoystickConfig      `json:"joystick"`
	SelfDriving    *SelfDrivingConfig   `json:"selfDriving"`
	SelfTracking   *SelfTrackingConfig  `json:"SelfTracking"`
	SpeechDriving  *SpeechDrivingConfig `json:"speechDriving"`
	SelfNav        *SelfNavConfig       `json:"selfnav"`
	BaiduAPIConfig *BaiduAPIConfig      `json:"baidu"`
	Volume         int                  `json:"volume"`
	Speed          uint32               `json:"speed"`
	Host           string               `json:"host"`
	VideoHost      string               `json:"videoHost"`
}

type L298NConfig struct {
	IN1Pin uint8 `json:"in1"`
	IN2Pin uint8 `json:"in2"`
	IN3Pin uint8 `json:"in3"`
	IN4Pin uint8 `json:"in4"`
	ENAPin uint8 `json:"ena"`
	ENBPin uint8 `json:"enb"`
}

type US100Config struct {
	Dev  string `json:"dev"`
	Baud int    `json:"baud"`
}

type GY25Config struct {
	Dev  string `json:"dev"`
	Baud int    `json:"baud"`
}

type LC12SConfig struct {
	Dev  string `json:"dev"`
	Baud int    `json:"baud"`
	CS   uint8  `json:"cs"`
}
type JoystickConfig struct {
	Enabled     bool         `json:"enabled"`
	LC12SConfig *LC12SConfig `json:"lc12s"`
}

type SelfDrivingConfig struct {
	Enabled bool `json:"enabled"`
}

type SelfTrackingConfig struct {
	Enabled   bool    `json:"enabled"`
	VideoHost string  `json:"videoHost"`
	LH        float64 `json:"lh"`
	LS        float64 `json:"ls"`
	LV        float64 `json:"lv"`
	HH        float64 `json:"hh"`
	HS        float64 `json:"hs"`
	HV        float64 `json:"hv"`
}

type SpeechDrivingConfig struct {
	Enabled bool `json:"enabled"`
}

type SelfNavConfig struct {
	Enabled       bool           `json:"enabled"`
	GPSConfig     *GPSConfig     `json:"gps"`
	TileMapConfig *TileMapConfig `json:"tileMap"`
}

type GPSConfig struct {
	Dev  string `json:"dev"`
	Baud int    `json:"baud"`
}

type TileMapConfig struct {
	MapFile  string    `json:"mapFile"`
	Box      *geo.Bbox `json:"bbox"`
	GridSize float64   `json:"gridSize"`
}

type BaiduAPIConfig struct {
	Speech *BaiduKey `json:"speech"`
	Image  *BaiduKey `json:"image"`
}

type BaiduKey struct {
	APIKey    string `json:"apiKey"`
	SecretKey string `json:"secretKey"`
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
