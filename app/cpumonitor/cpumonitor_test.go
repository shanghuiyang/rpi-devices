package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStart(t *testing.T) {
	monitor := cpuMonitor{}
	assert.NotNil(t, monitor)
}
