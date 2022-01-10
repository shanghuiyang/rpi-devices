package main

import (
	"fmt"
	"log"
	"time"

	"github.com/shanghuiyang/rpi-devices/dev"
	"github.com/shanghuiyang/rpi-devices/iot"
)

const (
	configJSON = "config.json"
)

type gardener struct {
	name       string
	workAtH    int
	workAtM    int
	workingSec int
	working    bool
	pump       dev.Pump
}

var (
	gardeners []*gardener
	buttom    dev.Button
	cloud     iot.Cloud
)

func main() {
	cfg, err := loadConfig(configJSON)
	if err != nil {
		log.Panicf("load config error: %v", err)
	}

	cloud = iot.NewNoop()
	if cfg.Iot.Enable {
		cloud = iot.NewOnenet(cfg.Iot.Onenet)
	}
	buttom = dev.NewButtonImp(cfg.Button)
	for _, g := range cfg.Gardeners {
		if !g.Enabled {
			continue
		}
		var h, m int
		if n, err := fmt.Sscanf(g.WorkAt, "%d:%d", &h, &m); n != 2 || err != nil {
			log.Panicf("parse watering time error: %v", err)
		}
		gardeners = append(gardeners, &gardener{
			name:       g.Name,
			workAtH:    h,
			workAtM:    m,
			workingSec: g.WorkingSec,
			pump:       dev.NewPumpImp(g.Pin),
		})
	}
	go timewater()
	go manwater()

	select {}
}

func timewater() {
	for {
		now := time.Now()
		h := now.Hour()
		m := now.Minute()
		for _, g := range gardeners {
			if g.workAtH == h && g.workAtM == m {
				go g.work()
			}
		}
		time.Sleep(time.Minute)
	}
}

func manwater() {
	for {
		if buttom.Pressed() {
			for _, g := range gardeners {
				go g.work()
				time.Sleep(time.Minute)
			}
			time.Sleep(100 * time.Millisecond)
		}
	}
}

func toCloud(g *gardener) {
	v := &iot.Value{
		Device: g.name,
		Value:  1,
	}
	if err := cloud.Push(v); err != nil {
		log.Printf("push to clould error: %v", err)
		return
	}

	time.Sleep(time.Duration(g.workingSec) * time.Second)
	v = &iot.Value{
		Device: g.name,
		Value:  0,
	}
	if err := cloud.Push(v); err != nil {
		log.Printf("push to clould error: %v", err)
		return
	}
	log.Printf("push to cloud successfully")
}

func (g *gardener) work() {
	if g.working {
		return
	}
	g.working = true
	g.pump.Run(g.workingSec)
	g.working = false
	toCloud(g)
	log.Printf("%v watered duration %v sec", g.name, g.workingSec)
}
