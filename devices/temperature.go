package devices

import (
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/shanghuiyang/pi/base"
	"github.com/shanghuiyang/pi/iotclouds"
)

const (
	logTagTemperature      = "temperature"
	lowTemperatureWarning  = 18
	highTemperatureWarning = 30
	temperatureInterval    = 1 * time.Minute
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

// Start ...
func (t *Temperature) Start() {
	log.Printf("[%v]start working", logTagTemperature)
	lastMail := time.Now()
	for {
		time.Sleep(temperatureInterval)
		c, err := t.GetTemperature()
		if err != nil {
			log.Printf("[%v]failed to get temperature, error: %v", logTagTemperature, err)
			continue
		}

		v := &iotclouds.IoTValue{
			DeviceName: TemperatureDevice,
			Value:      c,
		}
		iotclouds.IotCloud.Push(v)
		ChLedOp <- Blink

		if c >= 27.3 {
			ChRelayOp <- On
		} else {
			ChRelayOp <- Off
		}

		if c <= lowTemperatureWarning || c >= highTemperatureWarning {
			d := time.Now().Sub(lastMail)
			if d > 15*time.Minute {
				go SendEmail(c)
				lastMail = time.Now()
			}
		}
	}
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

// SendEmail ...
func SendEmail(temperatue float32) {
	subject := "Low Temperature Warning"
	if temperatue >= highTemperatureWarning {
		subject = "High Temperature Warning"
	}

	info := &base.EmailInfo{
		To:      base.GetEmailList(),
		Subject: subject,
		Body:    fmt.Sprintf("Current Temperature: %v 'C", temperatue),
	}
	base.SendEmail(info)
	ChLedOp <- Blink
}
