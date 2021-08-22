package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/shanghuiyang/rpi-devices/dev"
	"github.com/shanghuiyang/rpi-devices/iot"
	"github.com/shanghuiyang/rpi-devices/util"
)

const (
	dioPin  = 9
	rclkPin = 10
	sclkPin = 11
)

const (
	onenetToken = "your_onenet_token"
	onenetAPI   = "http://api.heclouds.com/devices/540381180/datapoints"
)

type data struct {
	name  string
	text  string
	value interface{}
}

type tempResponse struct {
	Temp     float32 `json:"temp"`
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
	dsp := dev.NewLedDisplay(dioPin, rclkPin, sclkPin)

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
			t, err := h.getTemp()
			if err != nil {
				log.Printf("[homeasst]failed to get temperature, error: %v", err)
				time.Sleep(5 * time.Second)
				return
			}
			log.Printf("[homeasst]temp: %v", t)
			d := &data{
				name:  "temp",
				text:  fmt.Sprintf("%.1f", t),
				value: t,
			}
			h.chDisplay <- d
			h.chCloud <- d
		}()

		go func() {
			pm25, err := h.getPM25()
			if err != nil {
				log.Printf("[homeasst]failed to get pm2.5, error: %v", err)
				time.Sleep(5 * time.Second)
				return
			}
			log.Printf("[homeasst]pm2.5: %v", pm25)

			d := &data{
				name:  "pm2.5",
				text:  fmt.Sprintf("%v", pm25),
				value: pm25,
			}
			h.chDisplay <- d
			h.chCloud <- d
			// h.chAlert <- v
		}()

		time.Sleep(60 * time.Second)
	}
}

func (h *homeAsst) display() {
	h.dsp.Open()
	opened := true
	cache := map[string]*data{}
	for {
		select {
		case d := <-h.chDisplay:
			cache[d.name] = d
		default:
			// do nothing, just use the latest temp
		}

		if h.dsp == nil {
			time.Sleep(30 * time.Second)
			continue
		}

		hour := time.Now().Hour()
		if hour >= 20 || hour < 8 {
			// turn off led display at 20:00-08:00
			if opened {
				h.dsp.Close()
				opened = false
			}
			time.Sleep(10 * time.Second)
			continue
		}

		if !opened {
			h.dsp.Open()
			opened = true
		}

		for _, d := range cache {
			h.dsp.Display(d.text)
			time.Sleep(5 * time.Second)
		}
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

func (h *homeAsst) getTemp() (float32, error) {
	resp, err := http.Get("http://localhost:8000/temp")
	if err != nil {
		return 0, fmt.Errorf("failed to get temp from sensers server, status: %v, err: %v", resp.Status, err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("failed to read resp body, err: %v", err)
	}

	var tempResp tempResponse
	if err := json.Unmarshal(body, &tempResp); err != nil {
		return 0, fmt.Errorf("failed to unmarshal resp, err: %v", err)
	}

	if tempResp.ErrorMsg != "" {
		return 0, fmt.Errorf("failed to get temp from sensers server, status: %v, err msg: %v", resp.Status, tempResp.ErrorMsg)
	}

	return tempResp.Temp, nil
}

func (h *homeAsst) getPM25() (uint16, error) {
	resp, err := http.Get("http://localhost:8000/pm25")
	if err != nil {
		return 0, fmt.Errorf("failed to get pm2.5 from sensers server, status: %v, err: %v", resp.Status, err)
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
