package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStart(t *testing.T) {
	dog := doordog{}
	assert.NotNil(t, dog)
}
