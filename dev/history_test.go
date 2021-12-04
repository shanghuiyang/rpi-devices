package dev

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewHistory(t *testing.T) {
	h := newHistory(1)
	assert.NotNil(t, h)
}

func TestHistoryAdd(t *testing.T) {
	testCases := []struct {
		desc    string
		element interface{}
	}{
		{
			desc:    "add uint8",
			element: uint8(1),
		},
		{
			desc:    "add uint16",
			element: uint16(1),
		},
		{
			desc:    "add uint32",
			element: uint16(1),
		},
		{
			desc:    "add uint32",
			element: uint32(1),
		},
		{
			desc:    "add uint64",
			element: uint64(1),
		},
		{
			desc:    "add int",
			element: int(1),
		},
		{
			desc:    "add int8",
			element: int8(1),
		},
		{
			desc:    "add int16",
			element: int16(1),
		},
		{
			desc:    "add int32",
			element: int32(1),
		},
		{
			desc:    "add int64",
			element: int64(1),
		},
		{
			desc:    "add float32",
			element: float32(1),
		},
		{
			desc:    "add float64",
			element: float64(1),
		},
		{
			desc:    "add string",
			element: "test",
		},
	}

	for _, test := range testCases {
		h := newHistory(10)
		assert.NotPanics(t, func() {
			h.Add(test.element)
		})
	}
}

func TestHistoryAvg(t *testing.T) {
	testCase := []struct {
		desc     string
		elements []interface{}
		noError  bool
		expected float64
	}{
		{
			desc:     "empty",
			elements: []interface{}{},
			noError:  false,
			expected: 0,
		},
		{
			desc:     "not full",
			elements: []interface{}{uint8(1)},
			noError:  true,
			expected: float64(1) / 1,
		},
		{
			desc:     "full",
			elements: []interface{}{int16(1), float32(2)},
			noError:  true,
			expected: float64(1+2) / 2,
		},
		{
			desc:     "full",
			elements: []interface{}{uint32(1), int32(2), float32(3)},
			noError:  true,
			expected: float64(2+3) / 2,
		},
		{
			desc:     "full",
			elements: []interface{}{uint8(1), uint16(2), uint32(3), uint64(4), int(5), int8(6), int16(7), int32(8), int64(9), float32(10), float64(11)},
			noError:  true,
			expected: float64(10+11) / 2,
		},
		{
			desc:     "string",
			elements: []interface{}{uint8(1), "some strings"},
			noError:  false,
			expected: 0,
		},
	}

	for _, test := range testCase {
		h := newHistory(2)
		assert.NotNil(t, h)

		for _, v := range test.elements {
			h.Add(v)
		}
		avg, err := h.Avg()
		if test.noError {
			assert.NoError(t, err)
			assert.Equal(t, test.expected, avg)
		} else {
			assert.Error(t, err)
		}
	}
}
