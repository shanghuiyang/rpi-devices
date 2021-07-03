package server

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCar(t *testing.T) {
	s := newService(nil, nil)
	assert.NotNil(t, s)
}
