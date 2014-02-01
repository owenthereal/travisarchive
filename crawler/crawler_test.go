package main

import (
	"testing"
	"time"

	"github.com/bmizerany/assert"
)

func TestCrawler_buildColName(t *testing.T) {
	form := "Jan 2, 2006 at 3:04pm (MST)"
	tt, _ := time.Parse(form, "Feb 3, 2013 at 7:54pm (UTC)")
	name := buildColName(&tt)

	assert.Equal(t, name, "builds_2013_02_03")
}

func TestCrawler_oneMinuteAgo(t *testing.T) {
	tt := oneMinuteAgo()
	now := time.Now()

	assert.T(t, now.After(tt))
}
