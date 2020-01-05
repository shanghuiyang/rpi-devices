package base

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReverse(t *testing.T) {
	testCases := []struct {
		desc     string
		text     string
		expected string
	}{
		{
			desc:     "normal case",
			text:     "abcd",
			expected: "dcba",
		},
		{
			desc:     "a case with space",
			text:     " 1 23 ",
			expected: " 32 1 ",
		},
		{
			desc:     "empty text",
			text:     "",
			expected: "",
		},
	}

	for _, test := range testCases {
		actual := Reverse(test.text)
		assert.Equal(t, test.expected, actual)
	}
}
