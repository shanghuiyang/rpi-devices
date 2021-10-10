package main

import (
	"bytes"
	"log"
	"net/http"

	"time"

	"image/color"

	sm "github.com/flopp/go-staticmaps"
	"github.com/golang/geo/s2"
	"github.com/shanghuiyang/rpi-devices/app/gpstracker/cache"
	"github.com/shanghuiyang/rpi-devices/dev"
	"github.com/shanghuiyang/rpi-devices/iot"
	"github.com/shanghuiyang/rpi-devices/util"
	"github.com/shanghuiyang/rpi-devices/util/geo"
)

const (
	zoomInBtnPin  = 18
	zoomOutBtnPin = 23

	minZoom    = 0
	maxZoom    = 17
	defautZoom = 17
	cacheDir   = ".cache/maptiles"

	streamerURL = ":8088/map"
	timeFormat  = "2006-01-02T15:04:05"
	onenetToken = "your_onenet_token"
	onenetAPI   = "http://api.heclouds.com/devices/540381180/datapoints"

	osmTiles             = "osm"
	googleSatelliteTiles = "google-satellite"
	bingSatelliteTiles   = "bing-satellite"
)

var defaultLoction = &geo.Point{
	Lat: 40.002369,
	Lon: 116.421977,
}

func main() {
	gps, err := dev.NewNeo6mGPS("/dev/ttyAMA0", 9600)
	// gps, err := dev.NewGPSSimulator("gps.csv")
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

	// cfg := &iot.Config{
	// 	Token: onenetToken,
	// 	API:   onenetAPI,
	// }
	// cloud := iot.NewOnenet(cfg)
	cloud := iot.NewNoop()

	streamer, err := util.NewStreamer(streamerURL)
	if err != nil {
		log.Printf("failed to create streamer, error: %v", err)
		return
	}
	zoomInBtn := dev.NewButtonImp(zoomInBtnPin)
	zoomOutBtn := dev.NewButtonImp(zoomOutBtnPin)
	t := &gpsTracker{
		gps:        gps,
		zoomInBtn:  zoomInBtn,
		zoomOutBtn: zoomOutBtn,
		logger:     logger,
		cloud:      cloud,
		streamer:   streamer,
		zoom:       defautZoom,
	}

	util.WaitQuit(t.close)
	t.start()
}

type gpsTracker struct {
	gps          dev.GPS
	cloud        iot.Cloud
	logger       util.Logger
	streamer     *util.Streamer
	zoomInBtn    dev.Button
	zoomOutBtn   dev.Button
	zoom         int
	tileProvider *sm.TileProvider
}

func (t *gpsTracker) start() {
	log.Printf("[gpstracker]start working")
	go t.detectZoomIn()
	go t.detectZoomOut()
	t.detectLoc()
}

func (t *gpsTracker) detectLoc() {
	c := cache.NewTileCache(cacheDir, 0777)
	t.tileProvider = sm.NewTileProviderOpenStreetMaps()
	m := util.NewMapRender()
	m.SetCache(c)
	m.SetTileProvider(t.tileProvider)
	m.SetSize(240, 240)

	lastPt := defaultLoction
	for {
		// time.Sleep(500 * time.Millisecond)
		pt, err := t.gps.Loc()
		if err != nil {
			log.Printf("[gpstracker]failed to get gps locations: %v", err)
			pt = lastPt
		}
		lastPt = pt

		t.logger.Printf("%v,%.6f,%.6f\n", time.Now().Format(timeFormat), pt.Lat, pt.Lon)

		v := &iot.Value{
			Device: "gps",
			Value:  pt,
		}
		go t.cloud.Push(v)

		marker := sm.NewMarker(
			s2.LatLngFromDegrees(pt.Lat, pt.Lon),
			color.RGBA{0xff, 0, 0, 0xff},
			12.0,
		)
		m.ClearMarker()
		m.AddMarker(marker)
		m.SetZoom(t.zoom)

		img, err := m.Render()
		if err != nil {
			log.Printf("[gpstracker]failed to render map: %v", err)
			util.DelayMs(100)
			continue
		}

		req, err := http.NewRequest("POST", "http://localhost:8080/display", bytes.NewBuffer(img))
		if err != nil {
			log.Printf("[gpstracker]failed to new http request: %v", err)
			util.DelayMs(100)
			continue
		}
		req.Header.Set("Content-Type", "application/json")
		client := &http.Client{
			Timeout: 1 * time.Second,
		}
		resp, err := client.Do(req)
		if err != nil {
			log.Printf("[gpstracker]failed to send http request: %v", err)
			util.DelayMs(200)
			continue
		}
		resp.Body.Close()
	}
}

func (t *gpsTracker) detectZoomIn() {
	n := 0
	for {
		if t.zoomInBtn.Pressed() {
			if n > 2 {
				// toggle tile type when keep pressing the button in 3s
				if t.tileProvider.Name == osmTiles {
					t.tileProvider.Name = bingSatelliteTiles
				} else {
					t.tileProvider.Name = osmTiles
				}
				log.Printf("[gpstracker]changed tile provider to: %v", t.tileProvider.Name)
				n = 0
				util.DelayMs(2000)
				continue
			}
			if n > 0 {
				n++
				util.DelayMs(1000)
				continue
			}

			n++
			t.zoomIn()
			log.Printf("[gpstracker]zoom in: z = %v", t.zoom)
			util.DelayMs(1000)
			continue
		}
		n = 0
		util.DelayMs(100)
	}
}

func (t *gpsTracker) detectZoomOut() {
	for {
		if t.zoomOutBtn.Pressed() {
			t.zoomOut()
			log.Printf("[gpstracker]zoom out: z = %v", t.zoom)
			util.DelayMs(1000)
			continue
		}
		util.DelayMs(100)
	}
}

func (t *gpsTracker) zoomIn() {
	if t.zoom >= maxZoom {
		return
	}
	t.zoom++
}

func (t *gpsTracker) zoomOut() {
	if t.zoom <= minZoom {
		return
	}
	t.zoom--
}

func (t *gpsTracker) close() {
	t.gps.Close()
	t.logger.Close()
}
