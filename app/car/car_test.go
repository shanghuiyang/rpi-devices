package main

import (
	"testing"

	dev "github.com/shanghuiyang/rpi-devices/devices"
	"github.com/stretchr/testify/assert"
)

func TestStart(t *testing.T) {
	car := dev.NewCar()
	assert.NotNil(t, car)
}
