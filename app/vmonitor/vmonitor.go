package main

import (
	"io/ioutil"
	"log"
	"net/http"

	"github.com/shanghuiyang/rpi-devices/base"
	"github.com/shanghuiyang/rpi-devices/dev"
	"github.com/stianeikeland/go-rpio"
)

const (
	pinSGH = 12
	pinSGV = 13
	pinLed = 4
	pinBzr = 10
)

var (
	vMonitor    *vmonitor
	pageContext []byte
)

func main() {
	if err := rpio.Open(); err != nil {
		log.Fatalf("failed to open rpio, error: %v", err)
		return
	}
	defer rpio.Close()

	hServo := dev.NewSG90(pinSGH)
	if hServo == nil {
		log.Printf("failed to new a sg90")
		return
	}

	vServo := dev.NewSG90(pinSGV)
	if vServo == nil {
		log.Printf("failed to new a sg90")
		return
	}

	led := dev.NewLed(pinLed)
	if led == nil {
		log.Printf("failed to new a led, will run the monitor without led")
	}

	bzr := dev.NewBuzzer(pinBzr)
	if bzr == nil {
		log.Printf("failed to new a buzzer, will run the monitor without buzzer")
	}

	vMonitor = newVMonitor(hServo, vServo, led, bzr)
	if vMonitor == nil {
		log.Printf("failed to new a vmonitor")
		return
	}

	if err := loadHomePage(); err != nil {
		log.Fatalf("failed to load home page, error: %v", err)
		return
	}

	log.Printf("video monitor server started")
	base.WaitQuit(func() {
		rpio.Close()
	})

	http.HandleFunc("/", handler)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err.Error())
	}
}

func loadHomePage() error {
	var err error
	pageContext, err = ioutil.ReadFile("vmonitor.html")
	if err != nil {
		return err
	}
	return nil
}

func handler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		homePageHandler(w, r)
	case "POST":
		operationHandler(w, r)
	}
}

func homePageHandler(w http.ResponseWriter, r *http.Request) {
	w.Write(pageContext)
	go vMonitor.led.Blink(5, 100)
}

func operationHandler(w http.ResponseWriter, r *http.Request) {
	op := r.FormValue("op")
	switch op {
	case "left":
		go vMonitor.left()
	case "right":
		go vMonitor.right()
	case "up":
		go vMonitor.up()
	case "down":
		go vMonitor.down()
	case "beep":
		go vMonitor.beep()
	default:
		log.Printf("invalid operator: %v", op)
	}
}

type vmonitor struct {
	hServo *dev.SG90
	vServo *dev.SG90
	led    *dev.Led
	buzzer *dev.Buzzer
	hAngle int
	vAngle int
}

func newVMonitor(hServo, vServo *dev.SG90, led *dev.Led, buzzer *dev.Buzzer) *vmonitor {
	v := &vmonitor{
		hServo: hServo,
		vServo: vServo,
		led:    led,
		buzzer: buzzer,
		hAngle: 0,
		vAngle: 0,
	}
	v.hServo.Roll(v.hAngle)
	v.vServo.Roll(v.vAngle)
	return v
}

func (v *vmonitor) left() {
	log.Printf("op: left")
	angle := v.hAngle - 15
	if angle < -90 {
		angle = -90
	}
	v.hAngle = angle
	log.Printf("servo: %v", angle)
	v.hServo.Roll(angle)
}

func (v *vmonitor) right() {
	log.Printf("op: right")
	angle := v.hAngle + 15
	if angle > 90 {
		angle = 90
	}
	v.hAngle = angle
	log.Printf("servo: %v", angle)
	v.hServo.Roll(angle)
}

func (v *vmonitor) up() {
	log.Printf("op: up")
	angle := v.vAngle + 15
	if angle > 90 {
		angle = 90
	}
	v.vAngle = angle
	log.Printf("servo: %v", angle)
	v.vServo.Roll(angle)
}

func (v *vmonitor) down() {
	log.Printf("op: down")
	angle := v.vAngle - 15
	if angle < -30 {
		return
	}
	v.vAngle = angle
	log.Printf("servo: %v", angle)
	v.vServo.Roll(angle)
}

func (v *vmonitor) beep() {
	log.Printf("op: beep")
	if v.buzzer == nil {
		return
	}
	v.buzzer.Beep(5, 100)
}
