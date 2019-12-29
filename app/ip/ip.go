package main

import (
	"fmt"
	"log"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/shanghuiyang/rpi-devices/base"
	"github.com/shanghuiyang/rpi-devices/iot"
)

const (
	maxRetry = 10
)

func main() {
	wsnCfg := &base.WsnConfig{
		Token: "your token",
		API:   "http://www.wsncloud.com/api/data/v1/numerical/insert",
	}
	cloud := iot.NewCloud(wsnCfg)
	if cloud == nil {
		log.Printf("failed to new OneNet iot cloud")
		return
	}

	n := 0
	for ; n < maxRetry; n++ {
		time.Sleep(10 * time.Second)
		ip, err := getIP()
		if err != nil {
			log.Printf("failed to get ip address, error: %v", err)
			log.Printf("retry %v...", n+1)
			continue
		}
		log.Printf("ip: %.6f", ip)

		v := &iot.Value{
			Device: "5e076e86e4b04a9a92a70f95",
			Value:  fmt.Sprintf("%.6f", ip),
		}
		if err := cloud.Push(v); err != nil {
			log.Printf("failed to push ip address to cloud, error: %v", err)
			log.Printf("retry %v...", n+1)
			continue
		}
		break
	}
	if n >= maxRetry {
		log.Printf("failed to get ip address")
		return
	}
	log.Printf("success")
}

func getIP() (float64, error) {
	cmd := exec.Command("hostname", "-I")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return 0, err
	}

	ip := strings.Trim(string(out), " \n")
	items := strings.Split(ip, ".")
	if len(items) != 4 {
		return 0, fmt.Errorf("failed to exec hostname")
	}
	ip1, err := strconv.Atoi(items[0])
	if err != nil {
		return 0, err
	}
	ip2, err := strconv.Atoi(items[1])
	if err != nil {
		return 0, err
	}
	ip3, err := strconv.Atoi(items[2])
	if err != nil {
		return 0, err
	}
	ip4, err := strconv.Atoi(items[3])
	if err != nil {
		return 0, err
	}
	result := float64(ip1)*1000 + float64(ip2) + float64(ip3)/1000 + float64(ip4)/1000000
	return result, nil
}
