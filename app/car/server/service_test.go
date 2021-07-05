package server

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCar(t *testing.T) {
	s, err := newService(&Config{})
	assert.NoError(t, err)
	assert.NotNil(t, s)
}
