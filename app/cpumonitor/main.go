package main

import (
	"errors"
	"log"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/shanghuiyang/rpi-devices/iot"
)

const (
	cpuInterval = 5 * time.Minute
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
	onenet := iot.NewOnenet(cfg)
	if onenet == nil {
		log.Printf("cpumonitor]failed to new OneNet iot cloud")
		return
	}
	monitor := &cpuMonitor{
		cloud: onenet,
	}
	monitor.start()
}

type cpuMonitor struct {
	cloud iot.Cloud
}

// Start ...
func (c *cpuMonitor) start() {
	log.Printf("[cpumonitor]cpu monitor start working")
	for {
		info, err := c.idle()
		if err != nil {
			log.Printf("[cpumonitor]failed to get cpu idle, error: %v", err)
			time.Sleep(30 * time.Second)
			continue
		}

		usage, err := c.usage(info)
		if err != nil {
			log.Printf("[cpumonitor]failed to get cpu usage, error: %v", err)
			time.Sleep(30 * time.Second)
			continue
		}

		v := &iot.Value{
			Device: "cpu",
			Value:  usage,
		}
		go c.cloud.Push(v)
		time.Sleep(cpuInterval)
	}
}

// Idle is to get idle cpu in %
// $ top -n 2 -d 1
// ---------------------------------------------------------------------------------
// top - 20:04:01 up 9 min,  2 users,  load average: 0.22, 0.22, 0.18
// Tasks:  72 total,   1 running,  71 sleeping,   0 stopped,   0 zombie
// %Cpu(s):  2.0 us,  2.0 sy,  0.0 ni, 96.0 id,  0.0 wa,  0.0 hi,  0.0 si,  0.0 st
// MiB Mem :    432.7 total,    330.8 free,     34.7 used,     67.2 buff/cache
// MiB Swap:    100.0 total,    100.0 free,      0.0 used.    347.1 avail Mem
// ---------------------------------------------------------------------------------
func (c *cpuMonitor) idle() (string, error) {
	cmd := exec.Command("top", "-b", "-n", "3", "-d", "3")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}

	return string(out), nil
}

func (c *cpuMonitor) usage(cpuinfo string) (float64, error) {
	lines := strings.Split(cpuinfo, "\n")
	var cpuline string
	for _, line := range lines {
		if strings.Contains(line, "Cpu") {
			cpuline = line
			break
		}
	}

	items := strings.Split(cpuline, ",")
	if len(items) != 8 {
		return 0, errors.New("invalid cup info")
	}
	id := items[3]
	id = strings.Trim(id, " ")
	id = strings.TrimRight(id, " id")
	idle, err := strconv.ParseFloat(id, 32)
	if err != nil {
		return 0, err
	}
	used := 100 - idle
	return used, nil
}
