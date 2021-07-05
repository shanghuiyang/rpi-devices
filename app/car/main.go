// build with tracking using open cv:
// $ go build -tags=gocv app/car/car.go

package main

import (
	"log"

	"github.com/shanghuiyang/rpi-devices/app/car/server"
)

const configJSON = "config.json"

func main() {
	cfg, err := server.LoadConfig(configJSON)
	if err != nil {
		log.Fatalf("failed to load config, error: %v", err)
		panic(err)
	}
	server.Start(cfg)
}
