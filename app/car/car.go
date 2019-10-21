package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	dev "github.com/shanghuiyang/rpi-devices/devices"
	"github.com/stianeikeland/go-rpio"
)

const (
	pinLed   = 4
	pinLight = 16
	pinIn1   = 17
	pinIn2   = 23
	pinIn3   = 27
	pinIn4   = 22
	pinENA   = 13
	pinENB   = 19
	pinBzr   = 10
	pinSG    = 18
)

var car *dev.Car

func main() {
	if err := rpio.Open(); err != nil {
		log.Fatalf("failed to open rpio, error: %v", err)
		return
	}
	defer rpio.Close()

	l298n := dev.NewL298N(pinIn1, pinIn2, pinIn3, pinIn4, pinENA, pinENB)
	if l298n == nil {
		log.Fatal("failed to new a L298N")
		return
	}
	buzzer := dev.NewBuzzer(pinBzr)
	if buzzer == nil {
		log.Printf("failed to new a buzzer, will build a car without horns")
	}

	led := dev.NewLed(pinLed)
	if led == nil {
		log.Printf("failed to new a led, will build a car without sign lights")
	}

	light := dev.NewLed(pinLight)
	if light == nil {
		log.Printf("failed to new a light, will build a car without lights")
	}

	sg := dev.NewSG90(pinSG)
	if sg == nil {
		log.Printf("failed to new a sg90, will build a camera without steering")
	}
	cam := dev.NewCamera(sg)
	if cam == nil {
		log.Printf("failed to new a sg90, will build a car without camera")
	}

	builder := dev.NewCarBuilder()
	car = builder.Engine(l298n).Horn(buzzer).Led(led).Light(light).Camera(cam).Build()
	if car == nil {
		log.Fatal("failed to new a car")
		return
	}
	car.Start()
	log.Printf("car server started")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		sig := <-quit
		log.Printf("received signal: " + sig.String() + ", stopping server")
		car.Stop()
		log.Printf("car server stoped")
		os.Exit(0)
	}()

	http.HandleFunc("/", carServer)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err.Error())
	}
}

func carServer(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		homePageHandler(w, r)
	case "POST":
		operationHandler(w, r)
	}
}

func homePageHandler(w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadFile("car.html")
	if err != nil {
		log.Printf("failed to read car.html")
		fmt.Fprintf(w, "failed to show home page")
	}
	w.Write(data)
}

func operationHandler(w http.ResponseWriter, r *http.Request) {
	op := r.FormValue("op")
	car.Do(dev.CarOp(op))
}
