package util

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

// Mode ...
type Mode string

// RpiModel ...
type RpiModel string

const (
	// DevMode ...
	DevMode = "dev"
	// PrdMode ...
	PrdMode = "prd"
)

const (
	// RpiUnknown ...
	RpiUnknown RpiModel = "Raspberry Pi X Model"
	// Rpi0 ...
	Rpi0 RpiModel = "Raspberry Pi Zero Model"
	// RpiA ...
	RpiA RpiModel = "Raspberry Pi A Model"
	// RpiB ...
	RpiB RpiModel = "Raspberry Pi B Model"
	// Rpi2 ...
	Rpi2 RpiModel = "Raspberry Pi 2 Model"
	// Rpi3 ...
	Rpi3 RpiModel = "Raspberry Pi 3 Model"
	// Rpi4 ...
	Rpi4 RpiModel = "Raspberry Pi 4 Model"
)

// Point is GPS point
type Point struct {
	Lat float32 `json:"lat"`
	Lon float32 `json:"lon"`
}

func (p *Point) String() string {
	return fmt.Sprintf("lat: %.6f, lon: %.6f", p.Lat, p.Lon)
}

// GetIP ...
func GetIP() string {
	cmd := exec.Command("hostname", "-I")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return ""
	}
	items := strings.Split(string(out), " ")
	if len(items) == 0 {
		return ""
	}
	return items[0]
}

// GetRpiModel ...
func GetRpiModel() RpiModel {
	cmd := exec.Command("cat", "/proc/device-tree/model")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return ""
	}
	s := string(out)
	if strings.Index(s, string(Rpi0)) >= 0 {
		return Rpi0
	}
	if strings.Index(s, string(RpiA)) >= 0 {
		return RpiA
	}
	if strings.Index(s, string(RpiB)) >= 0 {
		return RpiB
	}
	if strings.Index(s, string(Rpi2)) >= 0 {
		return Rpi2
	}
	if strings.Index(s, string(Rpi3)) >= 0 {
		return Rpi3
	}
	if strings.Index(s, string(Rpi4)) >= 0 {
		return Rpi4
	}
	return RpiUnknown
}

// Reverse reverses the string
func Reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

// WaitQuit ...
func WaitQuit(beforeQuitFunc func()) {
	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-c
		log.Printf("[base]received signal: %v, will quit\n", sig)
		beforeQuitFunc()
		os.Exit(0)
	}()
}

// DelayMs ...
func DelayMs(ms int) {
	time.Sleep(time.Duration(ms) * time.Millisecond)
}
