package main

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
	"time"

	"github.com/shanghuiyang/rpi-devices/iot"
)

const (
	memoryInterval = 10 * time.Minute
)

const (
	onenetToken = "your_onenet_token"
	onenetAPI   = "http://api.heclouds.com/devices/540381180/datapoints"
)

func main() {
	cfg := &iot.Config{
		Token: onenetToken,
		API:   onenetAPI,
	}
	cloud := iot.NewOnenet(cfg)
	if cloud == nil {
		log.Printf("[memmonitor]failed to new OneNet iot cloud")
		return
	}

	monitor := &memMonitor{
		cloud: cloud,
	}
	monitor.start()
}

type memMonitor struct {
	cloud iot.Cloud
}

func (m *memMonitor) start() {
	log.Printf("[memmonitor]start working")
	for {
		f, err := m.free()
		if err != nil {
			log.Printf("[memmonitor]failed to get free memory, error: %v", err)
			time.Sleep(30 * time.Second)
			continue
		}
		v := &iot.Value{
			Device: "memory",
			Value:  f,
		}
		go m.cloud.Push(v)
		time.Sleep(memoryInterval)
	}
}

// Free is to get free memory in MB
// $ free -m
// ---------------------------------------------------------------------------------
//             total        used        free      shared  buff/cache   available
// Mem:          432          50         258           3         123         328
// Swap:          99           0          99
// ---------------------------------------------------------------------------------
func (m *memMonitor) free() (float32, error) {
	cmd := exec.Command("free", "-m")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return 0, err
	}
	var v float32
	str := string(out)
	lines := strings.Split(str, "\n")
	if len(lines) < 3 {
		return 0, fmt.Errorf("failed to exec free")
	}
	items := strings.Split(lines[1], " ")
	if len(items) < 1 {
		return 0, fmt.Errorf("failed to parse")
	}
	if n, err := fmt.Sscanf(items[len(items)-1], "%f", &v); n != 1 || err != nil {
		return 0, fmt.Errorf("failed to parse")
	}
	return v, nil
}
