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
		Token: base.WsnToken,
		API:   base.WsnNumericalAPI,
	}
	cloud := iot.NewCloud(wsnCfg)
	if cloud == nil {
		log.Printf("[ip]failed to new OneNet iot cloud")
		return
	}

	n := 0
	for ; n < maxRetry; n++ {
		time.Sleep(10 * time.Second)
		ip, err := getIP()
		if err != nil {
			log.Printf("[ip]failed to get ip address, error: %v", err)
			log.Printf("[ip]retry %v...", n+1)
			continue
		}
		log.Printf("[ip]ip: %.6f", ip)

		v := &iot.Value{
			Device: "5e076e86e4b04a9a92a70f95",
			Value:  fmt.Sprintf("%.6f", ip),
		}
		if err := cloud.Push(v); err != nil {
			log.Printf("[ip]failed to push ip address to cloud, error: %v", err)
			log.Printf("[ip]retry %v...", n+1)
			continue
		}
		break
	}
	if n >= maxRetry {
		log.Printf("[ip]failed to get ip address")
		return
	}
	log.Printf("[ip]success")
}

func getIP() (float64, error) {
	cmd := exec.Command("hostname", "-I")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return 0, err
	}
	items := strings.Split(string(out), " ")
	if len(items) == 0 {
		return 0, fmt.Errorf("failed to exec hostname")
	}
	ips := strings.Split(items[0], ".")
	if len(ips) != 4 {
		return 0, fmt.Errorf("incorrect ip format")
	}
	ip1, err := strconv.Atoi(ips[0])
	if err != nil {
		return 0, err
	}
	ip2, err := strconv.Atoi(ips[1])
	if err != nil {
		return 0, err
	}
	ip3, err := strconv.Atoi(ips[2])
	if err != nil {
		return 0, err
	}
	ip4, err := strconv.Atoi(ips[3])
	if err != nil {
		return 0, err
	}
	result := float64(ip1)*1000 + float64(ip2) + float64(ip3)/1000 + float64(ip4)/1000000
	return result, nil
}
