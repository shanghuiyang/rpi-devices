package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/shanghuiyang/rpi-devices/base"
	"github.com/shanghuiyang/rpi-devices/dev"
	"github.com/shanghuiyang/rpi-devices/iot"
	"github.com/stianeikeland/go-rpio"
)

const (
	dioPin  = 9
	rclkPin = 10
	sclkPin = 11
)

type value struct {
	temp float32
	pm25 uint16
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
	dsp       *dev.LedDisplay
	cloud     iot.Cloud
	chDisplay chan *value // for disploying on oled
	chCloud   chan *value // for pushing to iot cloud
	chAlert   chan *value // for alerting
}

func main() {
	if err := rpio.Open(); err != nil {
		log.Fatalf("[homeasst]failed to open rpio, error: %v", err)
		return
	}
	defer rpio.Close()

	dsp := dev.NewLedDisplay(dioPin, rclkPin, sclkPin)

	onenetCfg := &base.OneNetConfig{
		Token: base.OneNetToken,
		API:   base.OneNetAPI,
	}
	cloud := iot.NewCloud(onenetCfg)

	asst := newHomeAsst(dsp, cloud)
	base.WaitQuit(func() {
		asst.stop()
		rpio.Close()
	})
	asst.start()
}

func newHomeAsst(dsp *dev.LedDisplay, cloud iot.Cloud) *homeAsst {
	return &homeAsst{
		dsp:       dsp,
		cloud:     cloud,
		chDisplay: make(chan *value, 4),
		chCloud:   make(chan *value, 4),
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
		t, err := h.getTemp()
		if err != nil {
			log.Printf("[homeasst]failed to get temperature, error: %v", err)
			time.Sleep(5 * time.Second)
			continue
		}
		log.Printf("[homeasst]temp: %v", t)

		pm25, err := h.getPM25()
		if err != nil {
			log.Printf("[homeasst]failed to get pm2.5, error: %v", err)
			time.Sleep(5 * time.Second)
			continue
		}
		log.Printf("[homeasst]pm2.5: %v", pm25)

		v := &value{
			temp: t,
			pm25: pm25,
		}
		h.chDisplay <- v
		h.chCloud <- v
		// h.chAlert <- v
		time.Sleep(60 * time.Second)
	}
}

func (h *homeAsst) display() {
	var v value
	h.dsp.Open()
	opened := true
	for {
		select {
		case vv := <-h.chDisplay:
			v = *vv
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

		tText := fmt.Sprintf("%.1f", v.temp)
		h.dsp.Display(tText)
		time.Sleep(5 * time.Second)

		pm25Text := fmt.Sprintf("%v", v.pm25)
		h.dsp.Display(pm25Text)
		time.Sleep(5 * time.Second)
	}
}

func (h *homeAsst) push() {
	for v := range h.chCloud {
		go func(v *value) {
			temp := &iot.Value{
				Device: "temp",
				Value:  v.temp,
			}
			if err := h.cloud.Push(temp); err != nil {
				log.Printf("[homeasst]push: failed to push temperature to cloud, error: %v", err)
			}

			pm25 := &iot.Value{
				Device: "pm2.5",
				Value:  v.pm25,
			}
			if err := h.cloud.Push(pm25); err != nil {
				log.Printf("[homeasst]push: failed to push pm2.5 to cloud, error: %v", err)
			}
		}(v)
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
		log.Printf("[homeasst]failed to get temp from sensers server, status: %v, err: %v", resp.Status, err)
		return 0, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("[homeasst]failed to read resp body, err: %v", err)
		return 0, err
	}

	var tempResp tempResponse
	if err := json.Unmarshal(body, &tempResp); err != nil {
		log.Printf("[homeasst]failed to unmarshal resp, err: %v", err)
		return 0, err
	}

	if tempResp.ErrorMsg != "" {
		log.Printf("[homeasst]failed to get temp from sensers server, status: %v, err msg: %v", resp.Status, tempResp.ErrorMsg)
		return 0, err
	}

	return tempResp.Temp, nil
}

func (h *homeAsst) getPM25() (uint16, error) {
	resp, err := http.Get("http://localhost:8000/pm25")
	if err != nil {
		log.Printf("[homeasst]failed to get pm2.5 from sensers server, status: %v, err: %v", resp.Status, err)
		return 0, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("[homeasst]failed to read resp body, err: %v", err)
		return 0, err
	}

	var pm25Resp pm25Response
	if err := json.Unmarshal(body, &pm25Resp); err != nil {
		log.Printf("[homeasst]failed to unmarshal resp, err: %v", err)
		return 0, err
	}

	if pm25Resp.ErrorMsg != "" {
		log.Printf("[homeasst]failed to get pm2.5 from sensers server, status: %v, err msg: %v", resp.Status, pm25Resp.ErrorMsg)
		return 0, err
	}

	return pm25Resp.PM25, nil
}
