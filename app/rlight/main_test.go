package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStart(t *testing.T) {
	light := rlight{}
	assert.NotNil(t, light)
}
