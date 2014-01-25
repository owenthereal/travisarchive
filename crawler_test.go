package main

import (
	"testing"
	"time"

	"github.com/bmizerany/assert"
)

func TestCrawler_buildColName(t *testing.T) {
	form := "Jan 2, 2006 at 3:04pm (MST)"
	tt, _ := time.Parse(form, "Feb 3, 2013 at 7:54pm (PST)")
	name := buildColName(&tt)

	assert.Equal(t, name, "builds_2013_02_04")
}
