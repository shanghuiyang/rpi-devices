package devices

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_GetTemperature(t *testing.T) {
	defer func(file string) {
		tempFile = file
	}(tempFile)

	tempFile = "./test/w1_slave"
	s := NewTemperature()
	assert.NotNil(t, s)

	v, err := s.GetTemperature()
	assert.NoError(t, err)
	assert.Equal(t, float32(28.625), v)
}
