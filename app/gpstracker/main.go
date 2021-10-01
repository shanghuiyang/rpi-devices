package main

import (
	"log"
	"time"

	"image/color"

	sm "github.com/flopp/go-staticmaps"
	"github.com/golang/geo/s2"
	"github.com/shanghuiyang/rpi-devices/dev"
	"github.com/shanghuiyang/rpi-devices/iot"
	"github.com/shanghuiyang/rpi-devices/util"
)

const (
	url         = ":8088/map"
	timeFormat  = "2006-01-02T15:04:05"
	onenetToken = "your_onenet_token"
	onenetAPI   = "http://api.heclouds.com/devices/540381180/datapoints"
)

func main() {
	gps, err := dev.NewNeo6mGPS("/dev/ttyAMA0", 9600)
	// gps, err := dev.NewGPSSimulator("./dev/test/gps.csv")
	if err != nil {
		log.Printf("[gpstracker]failed to new a gps device: %v", err)
		return
	}

	logfile := time.Now().Format(timeFormat) + ".csv"
	logger, err := util.NewGPSLogger(logfile)
	if err != nil {
		log.Printf("[gpstracker]failed to new gpslogger")
		return
	}
	// logger := util.NewNoopLogger()

	cfg := &iot.Config{
		Token: onenetToken,
		API:   onenetAPI,
	}
	cloud := iot.NewOnenet(cfg)
	// cloud := iot.NewNoop()

	streamer, err := util.NewStreamer(url)
	if err != nil {
		log.Printf("failed to create streamer, error: %v", err)
		return
	}
	t := &gpsTracker{
		gps:      gps,
		logger:   logger,
		cloud:    cloud,
		streamer: streamer,
	}

	util.WaitQuit(t.close)
	t.start()
}

type gpsTracker struct {
	gps      dev.GPS
	cloud    iot.Cloud
	logger   util.Logger
	streamer *util.Streamer
}

func (t *gpsTracker) start() {
	log.Printf("[gpstracker]start working")

	m := util.NewMapRender()
	m.SetSize(400, 400)
	m.SetZoom(16)

	for {
		time.Sleep(2 * time.Second)
		pt, err := t.gps.Loc()
		if err != nil {
			log.Printf("[gpstracker]failed to get gps locations: %v", err)
			continue
		}

		t.logger.Printf("%v,%.6f,%.6f\n", time.Now().Format(timeFormat), pt.Lat, pt.Lon)

		v := &iot.Value{
			Device: "gps",
			Value:  pt,
		}
		go t.cloud.Push(v)

		marker := sm.NewMarker(
			s2.LatLngFromDegrees(pt.Lat, pt.Lon),
			color.RGBA{0xff, 0, 0, 0xff},
			16.0,
		)
		m.ClearMarker()
		m.AddMarker(marker)
		img, err := m.Render()
		if err != nil {
			log.Printf("[gpstracker]failed to render map: %v", err)
			continue
		}
		if t.streamer != nil {
			t.streamer.Push(img)
		}
	}
}

func (t *gpsTracker) close() {
	t.gps.Close()
	t.logger.Close()
}
