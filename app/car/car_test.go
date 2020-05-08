package main

import (
	"testing"

	"github.com/shanghuiyang/rpi-devices/dev"
	"github.com/stretchr/testify/assert"
)

func TestStart(t *testing.T) {
	car := dev.NewCar()
	assert.NotNil(t, car)

	s := newCarServer(car)
	assert.NotNil(t, s)
}
