package main

import (
	"testing"

	"github.com/jakefau/rpi-devices/app/car/car"
	"github.com/stretchr/testify/assert"
)

func TestStart(t *testing.T) {
	car := car.New(&car.Config{})
	assert.NotNil(t, car)

	s := newServer(car)
	assert.NotNil(t, s)
}
