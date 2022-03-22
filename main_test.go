package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConvertDatetimeToUnixTimestamp(t *testing.T) {
	_, err := ConvertDatetimeToUnixTimestamp("2020-01-01")
	assert.Error(t, err)

	ts, err := ConvertDatetimeToUnixTimestamp("2022-01-01T00:00:00Z")
	assert.NoError(t, err)
	assert.Equal(t, int64(1640995200), ts)
}

func TestLargest(t *testing.T) {
	m := map[string]float64{
		"pool1": 0.1,
		"pool2": 0.2,
		"pool3": 0.3,
		"pool4": 0.1,
	}
	poolId, earningRate := Largest(m)
	assert.Equal(t, "pool3", poolId)
	assert.Equal(t, 0.3, earningRate)
}
