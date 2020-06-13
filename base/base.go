package base

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
)

// Mode ...
type Mode string

const (
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

// SendEmail ...
func SendEmail(info *EmailInfo) {
	chEmail <- info
}

// GetEmailList ...
func GetEmailList() []string {
	return emailList
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
