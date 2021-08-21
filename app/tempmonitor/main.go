package main

import (
	"fmt"
	"log"
	"os/exec"
	"time"

	"github.com/shanghuiyang/rpi-devices/dev"
	"github.com/shanghuiyang/rpi-devices/iot"
	"github.com/shanghuiyang/rpi-devices/util"
	"github.com/stianeikeland/go-rpio"
)

const (
	ledPin                 = 12
	lowTemperatureWarning  = 18
	highTemperatureWarning = 30
	intervalTime           = 1 * time.Minute
)

const (
	onenetToken = "your_onenet_token"
	onenetAPI   = "http://api.heclouds.com/devices/540381180/datapoints"
)

func main() {
	if err := rpio.Open(); err != nil {
		log.Fatalf("[tempmonitor]failed to open rpio, error: %v", err)
		return
	}
	defer rpio.Close()

	ds18b20 := dev.NewDS18B20()
	if ds18b20 == nil {
		log.Printf("[tempmonitor]failed to new temperature sensor")
		return
	}
	led := dev.NewLedImp(ledPin)
	if led == nil {
		log.Printf("[tempmonitor]failed to new led")
		return
	}

	cfg := &iot.Config{
		Token: onenetToken,
		API:   onenetAPI,
	}
	cloud := iot.NewOnenet(cfg)
	if cloud == nil {
		log.Printf("[tempmonitor]failed to new OneNet iot cloud")
		return
	}

	monitor := tempMonitor{
		thermometer: ds18b20,
		cloud:       cloud,
		led:         led,
	}

	util.WaitQuit(func() {
		rpio.Close()
	})

	monitor.start()
}

type tempMonitor struct {
	thermometer dev.Thermometer
	led         dev.Led
	cloud       iot.Cloud
}

func (m *tempMonitor) start() {
	for {
		time.Sleep(intervalTime)
		c, err := m.thermometer.Temperature()
		if err != nil {
			log.Printf("[tempmonitor]failed to get temperature, error: %v", err)
			continue
		}

		v := &iot.Value{
			Device: "temperature",
			Value:  c,
		}
		go m.cloud.Push(v)
		go m.led.Blink(5, 500)

		if c <= lowTemperatureWarning || c >= highTemperatureWarning {
			go m.notitfy(c)
		}
	}
}

func (m *tempMonitor) notitfy(temperatue float32) {
	_, err := exec.LookPath("mutt")
	if err != nil {
		log.Printf("[tempmonitor]need to install mutt for email notification")
		return
	}
	subject := "Low Temperature Warning"
	if temperatue >= highTemperatureWarning {
		subject = "High Temperature Warning"
	}
	subject = fmt.Sprintf("%v: %.2f C", subject, temperatue)
	cmd := exec.Command("mutt", "-s", subject, "youremail@xxx.com")
	if err := cmd.Run(); err != nil {
		log.Printf("[tempmonitor]failed to send email")
	}
}
