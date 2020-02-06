package main

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/shanghuiyang/rpi-devices/base"
	"github.com/shanghuiyang/rpi-devices/dev"
	"github.com/stianeikeland/go-rpio"
)

const (
	pinSG     = 18
	ipPattern = "((000.000.000.000))"
)

var (
	fan         *autoFan
	pageContext []byte
)

func main() {
	if err := rpio.Open(); err != nil {
		log.Fatalf("failed to open rpio, error: %v", err)
		return
	}
	defer rpio.Close()

	sg := dev.NewSG90(pinSG)
	if sg == nil {
		log.Printf("failed to new a sg90, will build a car without servo")
	}
	fan = newAuotFan(sg)

	if err := loadHomePage(); err != nil {
		log.Fatalf("failed to load home page, error: %v", err)
		return
	}

	log.Printf("fan server started")

	base.WaitQuit(func() {
		rpio.Close()
	})

	http.HandleFunc("/", fanServer)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err.Error())
	}
}

type autoFan struct {
	servo *dev.SG90
	state string // on or off
}

func newAuotFan(sg *dev.SG90) *autoFan {
	sg.Roll(-90)
	return &autoFan{
		servo: sg,
		state: "off",
	}
}

func loadHomePage() error {
	data, err := ioutil.ReadFile("car.html")
	if err != nil {
		return errors.New("internal error: failed to read car.html")
	}

	ip := base.GetIP()
	if ip == "" {
		return errors.New("internal error: failed to get ip")
	}

	rbuf := bytes.NewBuffer(data)
	wbuf := bytes.NewBuffer([]byte{})
	for {
		line, err := rbuf.ReadBytes('\n')
		if err == io.EOF {
			break
		}
		s := string(line)
		if strings.Index(s, ipPattern) >= 0 {
			s = strings.Replace(s, ipPattern, ip, 1)
		}
		wbuf.Write([]byte(s))
	}
	pageContext = wbuf.Bytes()
	return nil
}

func fanServer(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		homePageHandler(w, r)
	case "POST":
		operationHandler(w, r)
	}
}

func homePageHandler(w http.ResponseWriter, r *http.Request) {
	w.Write(pageContext)
}

func operationHandler(w http.ResponseWriter, r *http.Request) {
	op := r.FormValue("op")
	switch op {
	case "on":
		if fan.state != "on" {
			fan.servo.Roll(90)
			time.Sleep(500 * time.Millisecond)
			fan.servo.Roll(-90)
			time.Sleep(500 * time.Millisecond)
			fan.state = "on"
		}
	case "off":
		if fan.state != "off" {
			fan.state = "off"
		}
	default:
		log.Printf("invaild op: %v", op)
	}
}
