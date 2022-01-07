package main

import (
	"fmt"
	"log"
	"time"

	"github.com/shanghuiyang/rpi-devices/dev"
	"github.com/shanghuiyang/rpi-devices/iot"
	"github.com/shanghuiyang/rpi-devices/util"
)

const (
	configJSON = "config.json"
)

var (
	stepper dev.Stepper
	button  dev.Button
	cloud   iot.Cloud
	total   int
)

func main() {
	cfg, err := loadConfig(configJSON)
	if err != nil {
		log.Fatalf("failed to load config, error: %v", err)
		panic(err)
	}

	stepper = dev.NewBYJ2848(cfg.Stepper.In1, cfg.Stepper.In2, cfg.Stepper.In3, cfg.Stepper.In4)
	button = dev.NewButtonImp(cfg.Button)
	cloud = iot.NewNoop()
	if cfg.Iot.Enable {
		cloud = iot.NewOnenet(cfg.Iot.Onenet)
	}

	var h, m int
	if n, err := fmt.Sscanf(cfg.FeedAt, "%d:%d", &h, &m); n != 2 || err != nil {
		log.Fatalf("parse feed time error: %v", err)
	}
	log.Printf("feed at: %02d:%02d", h, m)

	go detectBtn()
	for {
		now := time.Now()
		if now.Hour() == h && now.Minute() == m {
			go feed()
		}
		time.Sleep(time.Minute)
	}
}

func detectBtn() {
	for {
		if button.Pressed() {
			go feed()
			util.DelayMs(1000)
		}
		util.DelayMs(100)
	}
}

func feed() {
	log.Printf("feeding")
	stepper.Roll(360)
	total++
	log.Printf("fed, total: %v", total)
	push()
}

func push() {
	v := &iot.Value{
		Device: "feeder",
		Value:  1,
	}
	if err := cloud.Push(v); err != nil {
		log.Printf("push to clould error: %v", err)
		return
	}

	time.Sleep(10 * time.Second)
	v = &iot.Value{
		Device: "feeder",
		Value:  0,
	}
	if err := cloud.Push(v); err != nil {
		log.Printf("push to clould error: %v", err)
		return
	}
	log.Printf("push to cloud successfully")
}
