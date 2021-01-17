package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStart(t *testing.T) {
	s := newTTSServer("app_key", "secret_key")
	assert.NotNil(t, s)
}
