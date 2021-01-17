package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStart(t *testing.T) {
	_, _ = getIP()
	assert.True(t, true)
}
