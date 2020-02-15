package main

import (
	"bytes"
	"fmt"
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
	pinSG           = 18
	statePattern    = "((state))"
	ipPattern       = "((000.000.000.000))"
	datetimePattern = "((yyyy-mm-dd hh:mm:ss))"
	datetimeFormat  = "2006-01-02 15:04:05"
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

func fanServer(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		homePageHandler(w, r)
	case "POST":
		operationHandler(w, r)
	}
}

func homePageHandler(w http.ResponseWriter, r *http.Request) {
	if len(pageContext) == 0 {
		var err error
		pageContext, err = ioutil.ReadFile("home.html")
		if err != nil {
			log.Printf("failed to read home.html")
			fmt.Fprintf(w, "internal error: failed to read home page")
			return
		}
	}

	ip := base.GetIP()
	if ip == "" {
		log.Printf("failed to get ip")
		fmt.Fprintf(w, "internal error: failed to get ip")
		return
	}

	wbuf := bytes.NewBuffer([]byte{})
	rbuf := bytes.NewBuffer(pageContext)
	for {
		line, err := rbuf.ReadBytes('\n')
		if err == io.EOF {
			break
		}
		s := string(line)
		switch {
		case strings.Index(s, ipPattern) >= 0:
			s = strings.Replace(s, ipPattern, ip, 1)
		case strings.Index(s, datetimePattern) >= 0:
			datetime := time.Now().Format(datetimeFormat)
			s = strings.Replace(s, datetimePattern, datetime, 1)
		case strings.Index(s, statePattern) >= 0:
			state := "unchecked"
			if fan.state == "on" {
				state = "checked"
			}
			s = strings.Replace(s, statePattern, state, 1)
		}
		wbuf.Write([]byte(s))
	}
	w.Write(wbuf.Bytes())
}

func operationHandler(w http.ResponseWriter, r *http.Request) {
	op := r.FormValue("op")
	log.Printf("op: %v", op)
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
