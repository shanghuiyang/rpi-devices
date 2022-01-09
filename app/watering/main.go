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
	cloud iot.Cloud
)

func main() {
	cfg, err := loadConfig(configJSON)
	if err != nil {
		log.Fatalf("failed to load config, error: %v", err)
		panic(err)
	}

	cloud = iot.NewNoop()
	if cfg.Iot.Enable {
		cloud = iot.NewOnenet(cfg.Iot.Onenet)
	}

	for _, c := range cfg.Pumps {
		go water(c)
	}

	select {}

}

func water(cfg *pumpConfig) {
	var h, m int
	if n, err := fmt.Sscanf(cfg.WateringAt, "%d:%d", &h, &m); n != 2 || err != nil {
		log.Panicf("parse watering time error: %v", err)
	}
	log.Printf("pump: %v, water at %02d:%02d, duration: %v sec", cfg.Name, h, m, cfg.WateringSec)

	var total int
	p := dev.NewPumpImp(cfg.Pin)

	// triggered by button
	go func() {
		if cfg.Button <= 0 {
			return
		}
		btn := dev.NewButtonImp(cfg.Button)
		for {
			if btn.Pressed() {
				go p.Run(cfg.WateringSec)
				go push(cfg.Name)
				total++
				log.Printf("pump %v watered duration %v sec, total: %v", cfg.Name, cfg.WateringSec, total)
				util.DelayMs(1000)
			}
			util.DelayMs(100)
		}
	}()

	// triggered by time
	for {
		now := time.Now()
		if now.Hour() == h && now.Minute() == m {
			go p.Run(cfg.WateringSec)
			go push(cfg.Name)
			log.Printf("pump %v watered duration %v sec, total: %v", cfg.Name, cfg.WateringSec, total)
			total++
		}
		time.Sleep(time.Minute)
	}
}

func push(name string) {
	v := &iot.Value{
		Device: name,
		Value:  1,
	}
	if err := cloud.Push(v); err != nil {
		log.Printf("push to clould error: %v", err)
		return
	}

	time.Sleep(10 * time.Second)
	v = &iot.Value{
		Device: name,
		Value:  0,
	}
	if err := cloud.Push(v); err != nil {
		log.Printf("push to clould error: %v", err)
		return
	}
	log.Printf("push to cloud successfully")
}
