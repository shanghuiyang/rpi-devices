package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/shanghuiyang/rpi-devices/devices"
)

var car *devices.Car

func main() {
	car = devices.NewCar()
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
		log.Printf(sig.String() + " received, stopping server")
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
	switch op {
	case "forward":
		car.Forward()
	case "backward":
		car.Backward()
	case "left":
		car.Left()
	case "right":
		car.Right()
	case "brake":
		car.Brake()
	default:
		log.Printf("invalid operation")
	}
}
