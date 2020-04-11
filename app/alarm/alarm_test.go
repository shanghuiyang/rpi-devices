package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStart(t *testing.T) {
	s := server{}
	assert.NotNil(t, s)
}
