package main

import (
	"testing"

	dev "github.com/shanghuiyang/rpi-devices/devices"
	"github.com/stretchr/testify/assert"
)

func TestStart(t *testing.T) {
	builder := dev.NewCarBuilder()
	assert.NotNil(t, builder)
	car = builder.Build()
	assert.NotNil(t, car)
}
