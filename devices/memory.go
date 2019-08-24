package devices

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
	"time"

	"github.com/shanghuiyang/pi/iotclouds"
)

const (
	logTagMemory   = "memory"
	memoryInterval = 10 * time.Minute
)

// Memory ...
type Memory struct {
}

// NewMemory ...
func NewMemory() *Memory {
	return &Memory{}
}

// Start ...
func (m *Memory) Start() {
	log.Printf("[%v]start working", logTagMemory)
	for {
		time.Sleep(memoryInterval)
		f, err := m.Free()
		if err != nil {
			log.Printf("[%v]failed to get free memory, error: %v", logTagMemory, err)
			continue
		}
		v := &iotclouds.IoTValue{
			DeviceName: MemoryDevice,
			Value:      f,
		}
		iotclouds.IotCloud.Push(v)
		ChLedOp <- Blink
	}
}

// Free is to get free memory in MB
// $ free -m
// ---------------------------------------------------------------------------------
//             total        used        free      shared  buff/cache   available
// Mem:          432          50         258           3         123         328
// Swap:          99           0          99
// ---------------------------------------------------------------------------------
func (m *Memory) Free() (float32, error) {
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
