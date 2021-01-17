package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStart(t *testing.T) {
	fan := autoFan{}
	assert.NotNil(t, fan)
}
