package base

import (
	"encoding/json"
	"io/ioutil"
)

// Config ...
type Config struct {
	Led       *LedConfig       `json:"led"`
	Relay     *RelayConfig     `json:"relay"`
	StepMotor *StepMotorConfig `json:"stepmotor"`
	Wsn       *WsnConfig       `json:"wsn"`
	OneNet    *OneNetConfig    `json:"onenet"`
	Email     *EmailConfig     `json:"email"`
	EmailTo   *EmailToConfig   `json:"emailto"`
}

// LedConfig ...
type LedConfig struct {
	Pin uint8 `json:"pin"`
}

// RelayConfig ...
type RelayConfig struct {
	Pin uint8 `json:"pin"`
}

// StepMotorConfig ...
type StepMotorConfig struct {
	In1 uint8 `json:"in1"`
	In2 uint8 `json:"in2"`
	In3 uint8 `json:"in3"`
	In4 uint8 `json:"in4"`
}

// WsnConfig ...
type WsnConfig struct {
	Token   string            `json:"token"`
	API     string            `json:"api"`
	Devices map[string]string `json:"devices"`
}

// OneNetConfig ...
type OneNetConfig struct {
	Token string `json:"token"`
	API   string `json:"api"`
}

// EmailConfig ...
type EmailConfig struct {
	SMTP     string `json:"smtp"`
	SMTPPort int    `json:"smtp_port"`
	POP      string `json:"pop"`
	POPPort  int    `json:"pop_port"`
	Address  string `json:"addr"`
	Password string `json:"password"`
}

// EmailToConfig ...
type EmailToConfig struct {
	List []string `json:"list"`
}

// LoadConfig ...
func LoadConfig() (*Config, error) {
	config := &Config{}
	data, err := ioutil.ReadFile("config.json")
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(data, config); err != nil {
		return nil, err
	}
	return config, nil
}
