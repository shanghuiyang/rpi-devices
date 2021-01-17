package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStart(t *testing.T) {
	h := homeAsst{}
	assert.NotNil(t, h)
}
