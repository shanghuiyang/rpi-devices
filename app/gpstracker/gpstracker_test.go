package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStart(t *testing.T) {
	gps := gpsTracker{}
	assert.NotNil(t, gps)
}
