package lbs

import (
	"fmt"
	"image"
	"log"
	"os"

	"image/color"
	"time"

	sm "github.com/flopp/go-staticmaps"
	"github.com/golang/geo/s2"
	"github.com/shanghuiyang/rpi-devices/app/gpstracker/tile"
	"github.com/shanghuiyang/rpi-devices/dev"
	"github.com/shanghuiyang/rpi-devices/iot"
	"github.com/shanghuiyang/rpi-devices/util"
	"github.com/shanghuiyang/rpi-devices/util/geo"
)

const (
	timeFormat = "2006-01-02T15:04:05"
)

var timer *time.Timer

type service struct {
	cfg             *Config
	gps             dev.GPS
	display         *dev.TFTDisplay
	cloud           iot.Cloud
	logger          util.Logger
	zoomInBtn       dev.Button
	zoomOutBtn      dev.Button
	tileProviders   map[string]*sm.TileProvider
	statusBarText   string
	curTileProvider *sm.TileProvider
	curZoom         int
	chPoint         chan *geo.Point
	chImage         chan image.Image
}

func newService(cfg *Config) (*service, error) {
	var gps dev.GPS
	gps, err := dev.NewHT1818GPS(cfg.GPS.Dev, cfg.GPS.Baud)
	if cfg.GPS.Simulator.Enable {
		latlons, e := loadPoints(cfg.GPS.Simulator.Source)
		if e != nil {
			log.Printf("[gpstracker]failed to load points from %v, error: %v", cfg.GPS.Simulator.Source, e)
			return nil, e
		}
		gps, err = dev.NewGPSSimulator(latlons)
	}

	if err != nil {
		log.Printf("[gpstracker]failed to new a gps device: %v", err)
		return nil, err
	}

	logfile := time.Now().Format(timeFormat) + ".csv"
	logger, err := util.NewGPSLogger(logfile)
	if err != nil {
		log.Printf("[gpstracker]failed to new gpslogger, error: %v", err)
		return nil, err
	}
	// logger := util.NewNoopLogger()

	display, err := dev.NewTFTDisplay(cfg.Display.Res, cfg.Display.Dc, cfg.Display.Blk, cfg.Display.Width, cfg.Display.Height)
	if err != nil {
		log.Printf("[gpstracker]failed to new display, error: %v", err)
		return nil, err
	}

	var cloud iot.Cloud = iot.NewNoop()
	if cfg.IOT.Enable {
		iotcfg := &iot.Config{
			Token: cfg.IOT.Onenet.Token,
			API:   cfg.IOT.Onenet.API,
		}
		cloud = iot.NewOnenet(iotcfg)
	}

	zoomInBtn := dev.NewButtonImp(cfg.ZoomInButtonPin)
	zoomOutBtn := dev.NewButtonImp(cfg.ZoomOutButtonPin)
	tileProviders := map[string]*sm.TileProvider{}
	for _, tileName := range cfg.Tile.TileProviders {
		tileProviders[tileName] = tile.NewLocalTileProvider(tileName)
	}

	return &service{
		cfg:             cfg,
		gps:             gps,
		display:         display,
		cloud:           cloud,
		logger:          logger,
		zoomInBtn:       zoomInBtn,
		zoomOutBtn:      zoomOutBtn,
		tileProviders:   tileProviders,
		curZoom:         cfg.Tile.DefaultZoom,
		curTileProvider: tileProviders[cfg.Tile.DefaultTileProvider],
		chPoint:         make(chan *geo.Point, 512),
		chImage:         make(chan image.Image, 512),
	}, nil
}

func (s *service) start() error {
	go s.detectZoomInBtn()
	go s.detectZoomOutBtn()
	go s.dispalyMap()
	go s.detectLocation()
	go s.renderMap()
	s.dispalyMap()
	return nil
}

func (s *service) detectLocation() {
	for {
		lat, lon, err := s.gps.Loc()
		if err != nil {
			log.Printf("failed to get gps locations: %v", err)
			continue
		}
		pt := &geo.Point{
			Lat: lat,
			Lon: lon,
		}

		s.chPoint <- pt
		v := &iot.Value{
			Device: "gps",
			Value:  pt,
		}
		go s.cloud.Push(v)
		s.logger.Printf("%v,%.6f,%.6f\n", time.Now().Format(timeFormat), pt.Lat, pt.Lon)

	}
}

// func (s *service) dispalyMap() {
// 	for img := range s.chImage {
// 		buf := &bytes.Buffer{}
// 		if err := png.Encode(buf, img); err != nil {
// 			log.Printf("failed to encode image, error: %v", err)
// 			continue
// 		}
// 		req, err := http.NewRequest("POST", "http://localhost:8080/display", buf)
// 		if err != nil {
// 			log.Printf("failed to new http request, error: %v", err)
// 			continue
// 		}
// 		req.Header.Set("Content-Type", "application/json")
// 		client := &http.Client{
// 			Timeout: 1 * time.Second,
// 		}
// 		resp, err := client.Do(req)
// 		if err != nil {
// 			log.Printf("failed to send http request, error: %v", err)
// 			continue
// 		}
// 		resp.Body.Close()
// 	}
// }

func (s *service) renderMap() {
	c := sm.NewTileCache(s.cfg.Tile.CacheDir, os.ModePerm)
	ctx := sm.NewContext()
	ctx.SetCache(c)
	ctx.SetOnline(s.cfg.Tile.Online)
	ctx.SetSize(s.cfg.Display.Width, s.cfg.Display.Height)

	updated := true
	lastZoom := s.curZoom
	lastProvider := s.curTileProvider
	lastPt := s.cfg.GPS.DefaultLoc
	curPt := s.cfg.GPS.DefaultLoc
	for {
		select {
		case curPt = <-s.chPoint:
			continue
		default:
		}

		if s.curZoom != lastZoom {
			lastZoom = s.curZoom
			updated = true
		}
		if s.curTileProvider != lastProvider {
			lastProvider = s.curTileProvider
			updated = true
		}
		if curPt.DistanceWith(lastPt) > 3 {
			lastPt = curPt
			updated = true
		}

		if !updated {
			util.DelayMs(1)
			continue
		}

		updated = false
		s.curTileProvider.Attribution = s.statusBarText
		marker := sm.NewMarker(
			s2.LatLngFromDegrees(curPt.Lat, curPt.Lon),
			color.RGBA{0xff, 0, 0, 0xff},
			12.0,
		)
		ctx.ClearObjects()
		ctx.AddObject(marker)
		ctx.SetZoom(s.curZoom)
		ctx.SetTileProvider(s.curTileProvider)

		img, err := ctx.Render()
		if err != nil {
			log.Printf("failed to render map: %v", err)
			util.DelayMs(100)
			continue
		}
		s.chImage <- img

	}
}

func (s *service) dispalyMap() {
	for img := range s.chImage {
		s.display.Display(img)
	}
}

func (s *service) toggleTileProvider() {
	provider := s.tileProviders[tile.OsmTile]
	if s.curTileProvider == provider {
		provider = s.tileProviders[tile.BingSatelliteTile]

	}
	s.curTileProvider = provider
	s.SetStatusBarText(fmt.Sprintf("Tile: %v", s.curTileProvider.Name))
	log.Printf("changed tile provider to: %v", provider.Name)
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
				util.DelayMs(400)
				continue
			}

			n++
			s.zoomIn()
			s.SetStatusBarText(fmt.Sprintf("Zoom: %v", s.curZoom))
			log.Printf("zoom in: z = %v", s.curZoom)
			util.DelayMs(500)
			continue
		}
		n = 0
		util.DelayMs(100)
	}
}

func (s *service) detectZoomOutBtn() {
	for {
		if s.zoomOutBtn.Pressed() {
			s.zoomOut()
			s.SetStatusBarText(fmt.Sprintf("Zoom: %v", s.curZoom))
			log.Printf("zoom out: z = %v", s.curZoom)
			util.DelayMs(500)
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

func (s *service) SetStatusBarText(text string) {
	if timer != nil {
		timer.Stop()
	}
	s.statusBarText = text

	// status bar will dispear after 5s
	timer = time.AfterFunc(5*time.Second, func() { s.statusBarText = "" })
}

func (s *service) Close() {
	s.gps.Close()
	s.display.Close()
	s.logger.Close()
}
