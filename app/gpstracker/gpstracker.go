package main

import (
	"log"
	"time"

	"github.com/shanghuiyang/rpi-devices/base"
	dev "github.com/shanghuiyang/rpi-devices/devices"
	"github.com/shanghuiyang/rpi-devices/iotclouds"
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
	cloud := iotclouds.New(oneNetCfg)
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
	cloud  iotclouds.IOTCloud
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
		v := &iotclouds.IoTValue{
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
