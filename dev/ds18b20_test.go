package dev

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_GetTemperature(t *testing.T) {
	defer func(file string) {
		tempFile = file
	}(tempFile)

	tempFile = "./test/w1_slave"
	d := NewDS18B20()
	assert.NotNil(t, d)

	v, err := d.Temperature()
	assert.NoError(t, err)
	assert.Equal(t, float64(28.625), v)
}
