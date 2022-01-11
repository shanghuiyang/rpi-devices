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
	workAt     string
	workingSec int
	working    bool
	relay      dev.Relay
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
		gardeners = append(gardeners, &gardener{
			name:       g.Name,
			workAt:     g.WorkAt,
			workingSec: g.WorkingSec,
			relay:      dev.NewRelayImp([]uint8{g.Relay}),
		})
	}

	go timewater()
	go manwater()

	select {}
}

func timewater() {
	for {
		now := time.Now()
		hm := fmt.Sprintf("%d:%d", now.Hour(), now.Minute())
		log.Printf("now: %v", hm)
		for _, g := range gardeners {
			if g.workAt == hm {
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
		log.Printf("push %v to clould error: %v", g.name, err)
		return
	}

	time.Sleep(time.Duration(g.workingSec) * time.Second)
	v = &iot.Value{
		Device: g.name,
		Value:  0,
	}
	if err := cloud.Push(v); err != nil {
		log.Printf("push %v to clould error: %v", g.name, err)
		return
	}
	log.Printf("push %v to cloud successfully", g.name)
}

func (g *gardener) work() {
	if g.working {
		return
	}
	log.Printf("%v is watering", g.name)
	g.working = true
	g.relay.On(0)
	time.Sleep(time.Duration(g.workingSec) * time.Second)
	g.relay.Off(0)
	g.working = false
	log.Printf("%v watered duration %v sec", g.name, g.workingSec)
	go toCloud(g)
}
