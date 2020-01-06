package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStart(t *testing.T) {
	m := ch2oMonitor{}
	assert.NotNil(t, m)
}
