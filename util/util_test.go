package util

import (
	"testing"

	"github.com/bmizerany/assert"
)

func TestParseBuildTime(t *testing.T) {
	tt, _ := ParseBuildTime("builds_2014_01_17")
	assert.Equal(t, "2014-01-17", tt.Format("2006-01-02"))

	_, err := ParseBuildTime("invalid")
	assert.NotEqual(t, nil, err)
}
