package util

import (
	"fmt"
	"log"
	"math"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

// Mode ...
type Mode string

const (
	// Epsilon ...
	Epsilon = 1e-9

	// DevMode ...
	DevMode = "dev"
	// PrdMode ...
	PrdMode = "prd"
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
func DelayMs(d time.Duration) {
	time.Sleep(d * time.Millisecond)
}

// IncludedAngle caculates the included angle betweet fromAngle to toAngle
func IncludedAngle(from, to float64) float64 {
	if from*to > 0 {
		return math.Abs(from - to)
	}

	d := math.Abs(from) + math.Abs(to)
	if d <= 180 {
		return d
	}
	return 360 - d
}

func AlmostEqual(a, b float64) bool {
	return math.Abs(a-b) < Epsilon
}
