package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStart(t *testing.T) {
	monitor := cpuMonitor{}
	assert.NotNil(t, monitor)
}

func TestUsage(t *testing.T) {
	cupinfo := `
		top - 20:04:01 up 9 min,  2 users,  load average: 0.22, 0.22, 0.18
		Tasks:  72 total,   1 running,  71 sleeping,   0 stopped,   0 zombie
		%Cpu(s):  2.0 us,  2.0 sy,  0.0 ni, 96.0 id,  0.0 wa,  0.0 hi,  0.0 si,  0.0 st
		MiB Mem :    432.7 total,    330.8 free,     34.7 used,     67.2 buff/cache
		MiB Swap:    100.0 total,    100.0 free,      0.0 used.    347.1 avail Mem
	`

	monitor := cpuMonitor{}
	assert.NotNil(t, monitor)

	usage, err := monitor.usage(cupinfo)
	assert.NoError(t, err)
	assert.InDelta(t, 4.0, usage, 1e-9)
}
