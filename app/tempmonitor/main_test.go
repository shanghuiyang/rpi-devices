package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStart(t *testing.T) {
	monitor := tempMonitor{}
	assert.NotNil(t, monitor)
}
