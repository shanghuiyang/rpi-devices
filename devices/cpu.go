package devices

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
	"time"

	"github.com/shanghuiyang/rpi-devices/iotclouds"
)

const (
	logTagCPU   = "cpu"
	cpuInterval = 5 * time.Minute
)

// CPU ...
type CPU struct {
}

// NewCPU ...
func NewCPU() *CPU {
	return &CPU{}
}

// Start ...
func (c *CPU) Start() {
	log.Printf("[%v]start working", logTagCPU)
	for {
		time.Sleep(cpuInterval)
		f, err := c.Idle()
		if err != nil {
			log.Printf("[%v]failed to get cpu idle, error: %v", logTagCPU, err)
			continue
		}

		v := &iotclouds.IoTValue{
			DeviceName: CPUDevice,
			Value:      f,
		}
		iotclouds.IotCloud.Push(v)
		ChLedOp <- Blink

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
func (c *CPU) Idle() (float32, error) {
	cmd := exec.Command("top", "-b", "-n", "3", "-d", "3")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return -1, err
	}
	str := string(out)
	lines := strings.Split(str, "\n")
	var cpuline string
	for _, line := range lines {
		if strings.Contains(line, "Cpu") {
			cpuline = line
		}
	}
	var cpu string
	items := strings.Split(cpuline, " ")
	for i, item := range items {
		if item == "id," && i > 0 {
			cpu = items[i-1]
		}
	}
	var v float32
	if n, err := fmt.Sscanf(cpu, "%f", &v); n != 1 || err != nil {
		return 0, fmt.Errorf("failed to parse")
	}
	return v, nil
}
