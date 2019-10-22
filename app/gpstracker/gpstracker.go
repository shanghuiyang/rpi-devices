package main

import (
	"log"
	"time"

	"github.com/shanghuiyang/rpi-devices/base"
	"github.com/shanghuiyang/rpi-devices/dev"
	"github.com/shanghuiyang/rpi-devices/iot"
)

func main() {
	gps := dev.NewGPS()
	if gps == nil {
		log.Printf("failed to new a gps device")
		return
	}
	logger := dev.NewGPSLogger()
	if logger == nil {
		log.Printf("failed to new a tracker")
		return
	}
	oneNetCfg := &base.OneNetConfig{
		Token: "your token",
		API:   "http://api.heclouds.com/devices/540381180/datapoints",
	}
	cloud := iot.NewCloud(oneNetCfg)
	if cloud == nil {
		log.Printf("failed to new OneNet iot cloud")
		return
	}
	t := &gpsTracker{
		gps:    gps,
		logger: logger,
		cloud:  cloud,
	}
	defer t.close()

	t.start()
}

type gpsTracker struct {
	gps    *dev.GPS
	logger *dev.GPSLogger
	cloud  iot.Cloud
}

func (t *gpsTracker) start() {
	log.Printf("gps tracker start working")
	for {
		time.Sleep(5 * time.Second)
		// pt, err := g.MockLocFromCSV()
		pt, err := t.gps.Loc()
		if err != nil {
			log.Printf("failed to get gps locations: %v", err)
			continue
		}
		t.logger.AddPoint(pt)
		v := &iot.Value{
			DeviceName: "gps",
			Value:      pt,
		}
		go t.cloud.Push(v)
	}
}

func (t *gpsTracker) close() {
	t.gps.Close()
	t.logger.Close()
}
