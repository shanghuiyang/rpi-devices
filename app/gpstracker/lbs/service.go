package lbs

import (
	"bytes"
	"log"
	"os"

	"image/color"
	"net/http"
	"time"

	"github.com/golang/geo/s2"
	sm "github.com/shanghuiyang/go-staticmaps"
	"github.com/shanghuiyang/rpi-devices/app/gpstracker/tile"
	"github.com/shanghuiyang/rpi-devices/dev"
	"github.com/shanghuiyang/rpi-devices/iot"
	"github.com/shanghuiyang/rpi-devices/util"
)

const (
	timeFormat = "2006-01-02T15:04:05"
)

type service struct {
	cfg             *Config
	gps             dev.GPS
	cloud           iot.Cloud
	logger          util.Logger
	zoomInBtn       dev.Button
	zoomOutBtn      dev.Button
	tileProviders   map[string]*sm.TileProvider
	curTileProvider *sm.TileProvider
	curZoom         int
	chImage         chan []byte
}

func newService(cfg *Config) (*service, error) {
	gps, err := dev.NewNeo6mGPS(cfg.GPS.Dev, cfg.GPS.Baud)
	// gps, err := dev.NewGPSSimulator("gps.csv")
	if err != nil {
		log.Printf("[gpstracker]failed to new a gps device: %v", err)
		return nil, err
	}

	logfile := time.Now().Format(timeFormat) + ".csv"
	logger, err := util.NewGPSLogger(logfile)
	if err != nil {
		log.Printf("[gpstracker]failed to new gpslogger")
		return nil, err
	}
	// logger := util.NewNoopLogger()

	// cfg := &iot.Config{
	// 	Token: onenetToken,
	// 	API:   onenetAPI,
	// }
	// cloud := iot.NewOnenet(cfg)
	cloud := iot.NewNoop()
	zoomInBtn := dev.NewButtonImp(cfg.ZoomInButtonPin)
	zoomOutBtn := dev.NewButtonImp(cfg.ZoomOutButtonPin)
	tileProviders := map[string]*sm.TileProvider{}
	for _, tileName := range cfg.Tile.TileProviders {
		tileProviders[tileName] = tile.NewLocalTileProvider(tileName)
	}

	return &service{
		cfg:             cfg,
		gps:             gps,
		cloud:           cloud,
		logger:          logger,
		zoomInBtn:       zoomInBtn,
		zoomOutBtn:      zoomOutBtn,
		tileProviders:   tileProviders,
		curZoom:         cfg.Tile.DefaultZoom,
		curTileProvider: tileProviders[cfg.Tile.DefaultTileProvider],
		chImage:         make(chan []byte),
	}, nil
}

func (s *service) start() error {
	go s.detectZoomInBtn()
	go s.detectZoomOutBtn()
	go s.dispalyMap()
	s.detectLocation()
	return nil
}

func (s *service) detectLocation() {
	c := sm.NewTileCache(s.cfg.Tile.CacheDir, os.ModePerm)
	// f := sm.NewTileFetcher(s.curTileProvider, c, s.cfg.Online)
	r := tile.NewRender()
	// r.SetTileFetcher(f)
	r.SetCache(c)
	r.SetOnline(s.cfg.Online)
	r.SetSize(240, 240)

	lastPt := s.cfg.DefaultLocation
	for {
		// time.Sleep(500 * time.Millisecond)
		pt, err := s.gps.Loc()
		if err != nil {
			log.Printf("[locator]failed to get gps locations: %v", err)
			pt = lastPt
		}
		lastPt = pt

		s.logger.Printf("%v,%.6f,%.6f\n", time.Now().Format(timeFormat), pt.Lat, pt.Lon)

		v := &iot.Value{
			Device: "gps",
			Value:  pt,
		}
		go s.cloud.Push(v)

		marker := sm.NewMarker(
			s2.LatLngFromDegrees(pt.Lat, pt.Lon),
			color.RGBA{0xff, 0, 0, 0xff},
			12.0,
		)
		r.ClearMarker()
		r.AddMarker(marker)
		r.SetZoom(s.curZoom)
		r.SetTileProvider(s.curTileProvider)

		img, err := r.Render()
		if err != nil {
			log.Printf("[locator]failed to render map: %v", err)
			util.DelayMs(100)
			continue
		}
		s.chImage <- img
	}
}

func (s *service) detectZoomInBtn() {
	n := 0
	for {
		if s.zoomInBtn.Pressed() {
			if n > 2 {
				// toggle tile type when keep pressing the button in 3s
				s.toggleTileProvider()
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
			s.zoomIn()
			log.Printf("[locator]zoom in: z = %v", s.curZoom)
			util.DelayMs(1000)
			continue
		}
		n = 0
		util.DelayMs(100)
	}
}

func (s *service) dispalyMap() {
	for img := range s.chImage {
		req, err := http.NewRequest("POST", "http://localhost:8080/display", bytes.NewBuffer(img))
		if err != nil {
			log.Printf("[locator]failed to new http request, error: %v", err)
			continue
		}
		req.Header.Set("Content-Type", "application/json")
		client := &http.Client{
			Timeout: 1 * time.Second,
		}
		resp, err := client.Do(req)
		if err != nil {
			log.Printf("[locator]failed to send http request, error: %v", err)
			continue
		}
		resp.Body.Close()
	}
}

func (s *service) toggleTileProvider() {
	provider := s.tileProviders[tile.OsmTile]
	if s.curTileProvider == provider {
		provider = s.tileProviders[tile.BingSatelliteTile]

	}
	s.curTileProvider = provider
	log.Printf("[locator]changed tile provider to: %v", provider.Name)
}

func (s *service) detectZoomOutBtn() {
	for {
		if s.zoomOutBtn.Pressed() {
			s.zoomOut()
			log.Printf("[locator]zoom out: z = %v", s.curZoom)
			util.DelayMs(1000)
			continue
		}
		util.DelayMs(100)
	}
}

func (s *service) zoomIn() {
	if s.curZoom >= s.cfg.Tile.MaxZoom {
		return
	}
	s.curZoom++
}

func (s *service) zoomOut() {
	if s.curZoom <= s.cfg.Tile.MinZoom {
		return
	}
	s.curZoom--
}

func (s *service) Close() {
	s.gps.Close()
	s.logger.Close()
}
