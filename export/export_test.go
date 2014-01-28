package main

import (
	"testing"

	"github.com/bmizerany/assert"
)

func TestParseDate(t *testing.T) {
	tt, _ := parseDate("builds_2014_01_17")
	assert.Equal(t, "2014-01-17", tt.Format("2006-01-02"))

	_, err := parseDate("invalid")
	assert.NotEqual(t, nil, err)
}
