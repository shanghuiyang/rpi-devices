package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"time"

	"github.com/shanghuiyang/rpi-devices/dev"
	"github.com/shanghuiyang/rpi-devices/iot"
	"github.com/shanghuiyang/rpi-devices/util"
)

const (
	onenetToken = "your_onenet_token"
	onenetAPI   = "http://api.heclouds.com/devices/540381180/datapoints"
)

var (
	celsiusChar        = []byte{0xDF}
	celsiusStr  string = string(celsiusChar[:])
)

type data struct {
	name        string
	value       interface{}
	displayText string
	displayX    int
	displayY    int
}

type temphumiResponse struct {
	Temp     float64 `json:"temp"`
	Humi     float64 `json:"humi"`
	ErrorMsg string  `json:"error_msg"`
}

type pm25Response struct {
	PM25     uint16 `json:"pm25"`
	ErrorMsg string `json:"error_msg"`
}

type homeAsst struct {
	dsp       dev.Display
	cloud     iot.Cloud
	chDisplay chan *data // for disploying on oled
	chCloud   chan *data // for pushing to iot cloud
	// chAlert   chan *data // for alerting
}

func main() {
	dsp, err := dev.NewLcdDisplay(16, 2)
	if err != nil {
		log.Printf("failed to new lcd display, error: %v", err)
		return
	}

	cfg := &iot.Config{
		Token: onenetToken,
		API:   onenetAPI,
	}
	onenet := iot.NewOnenet(cfg)

	asst := newHomeAsst(dsp, onenet)
	util.WaitQuit(func() {
		asst.stop()
	})
	asst.start()
}

func newHomeAsst(dsp dev.Display, cloud iot.Cloud) *homeAsst {
	return &homeAsst{
		dsp:       dsp,
		cloud:     cloud,
		chDisplay: make(chan *data, 4),
		chCloud:   make(chan *data, 4),
		// chAlert:   make(chan *value, 4),
	}
}

func (h *homeAsst) start() {
	go h.display()
	go h.push()
	// go h.alert()
	h.getData()
}

func (h *homeAsst) getData() {
	for {
		go func() {
			temp, humi, err := h.getTempHumi()
			if err != nil {
				log.Printf("[homeasst]failed to get temp and humi, error: %v", err)
				time.Sleep(5 * time.Second)
				return
			}
			log.Printf("[homeasst]temp: %v, humi: %v", temp, humi)
			d := &data{
				name:        "temp",
				value:       math.Round(temp*100) / 100,
				displayText: fmt.Sprintf("%5.1f%sC", temp, celsiusStr),
				displayX:    0,
				displayY:    0,
			}
			h.chDisplay <- d
			h.chCloud <- d

			d = &data{
				name:        "humi",
				value:       math.Round(humi*100) / 100,
				displayText: fmt.Sprintf("%5.1f%%", humi),
				displayX:    10,
				displayY:    0,
			}
			h.chDisplay <- d
			h.chCloud <- d
		}()

		// go func() {
		// 	pm25, err := h.getPM25()
		// 	if err != nil {
		// 		log.Printf("[homeasst]failed to get pm2.5, error: %v", err)
		// 		time.Sleep(5 * time.Second)
		// 		return
		// 	}
		// 	log.Printf("[homeasst]pm2.5: %v", pm25)

		// 	d := &data{
		// 		name:        "pm2.5",
		// 		value:       pm25,
		// 		displayText: fmt.Sprintf("Air:%3d", pm25),
		// 		displayX:    8,
		// 		displayY:    0,
		// 	}
		// 	h.chDisplay <- d
		// 	h.chCloud <- d
		// 	// h.chAlert <- v
		// }()

		time.Sleep(60 * time.Second)
	}
}

func (h *homeAsst) display() {
	h.dsp.On()
	backlight := true

	go func() {
		for {
			d := &data{
				name:        "time",
				displayText: time.Now().Format("15:04:05"),
				displayX:    4,
				displayY:    1,
			}
			h.chDisplay <- d
			time.Sleep(1 * time.Second)
		}
	}()

	cache := map[string]*data{}
	for {
		select {
		case d := <-h.chDisplay:
			cache[d.name] = d
			continue // continue to update the data in cache
		default:
			// do nothing, just use the latest temp
		}

		hour := time.Now().Hour()
		if hour >= 20 || hour < 8 {
			// turn off backlight at 20:00-08:00
			if backlight {
				h.dsp.Off()
				backlight = false
			}
		} else if !backlight {
			h.dsp.On()
			backlight = true
		}

		for _, d := range cache {
			h.dsp.Text(d.displayText, d.displayX, d.displayY)
		}
		time.Sleep(1 * time.Second)
	}
}

func (h *homeAsst) push() {
	for d := range h.chCloud {
		go func(d *data) {
			v := &iot.Value{
				Device: d.name,
				Value:  d.value,
			}
			if err := h.cloud.Push(v); err != nil {
				log.Printf("[homeasst]failed to push %v to cloud, error: %v", d.name, err)
				return
			}
		}(d)
	}
}

// func (h *homeAsst) alert() {
// 	var temp, humi float32 = -999, -999
// 	for {
// 		select {
// 		case v := <-h.chAlert:
// 			temp, humi = v.temp, v.humi
// 		default:
// 			// do nothing
// 		}

// 		if (temp > 0 && temp < 15) || humi > 70 {
// 			h.led.Blink(1, 1000)
// 			continue
// 		}
// 		time.Sleep(1 * time.Second)
// 	}
// }

func (h *homeAsst) stop() {
	h.dsp.Close()
}

func (h *homeAsst) getTempHumi() (temp, humi float64, err error) {
	resp, err := http.Get("http://localhost:8000/temphumi")
	if err != nil {
		return 0, 0, fmt.Errorf("failed to get temp from sensers server, err: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to read resp body, err: %v", err)
	}

	var result temphumiResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return 0, 0, fmt.Errorf("failed to unmarshal resp, err: %v", err)
	}

	if result.ErrorMsg != "" {
		return 0, 0, fmt.Errorf("failed to get temp from sensers server, status: %v, err msg: %v", resp.Status, result.ErrorMsg)
	}

	return result.Temp, result.Humi, nil
}

func (h *homeAsst) getPM25() (uint16, error) {
	resp, err := http.Get("http://localhost:8000/pm25")
	if err != nil {
		return 0, fmt.Errorf("failed to get pm2.5 from sensers server, err: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("failed to read resp body, err: %v", err)
	}

	var pm25Resp pm25Response
	if err := json.Unmarshal(body, &pm25Resp); err != nil {
		return 0, fmt.Errorf("failed to unmarshal resp, err: %v", err)
	}

	if pm25Resp.ErrorMsg != "" {
		return 0, fmt.Errorf("failed to get pm2.5 from sensers server, status: %v, err msg: %v", resp.Status, pm25Resp.ErrorMsg)
	}

	return pm25Resp.PM25, nil
}
