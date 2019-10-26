package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStart(t *testing.T) {
	light := autoLight{}
	assert.NotNil(t, light)
}
