package main

import (
	"net/url"
	"testing"

	"github.com/bmizerany/assert"
)

func TestMongoRxport_buildArgs(t *testing.T) {
	u := "mongodb://user:pass@zach.mongohq.com:10081/travisarchive"
	uu, _ := url.Parse(u)
	col := "builds_123"

	cmd := &MongoExport{ColName: col}
	args := cmd.buildArgs(uu, "file")

	assert.Equal(t, "--port", args[0])
	assert.Equal(t, "10081", args[1])
	assert.Equal(t, "-h", args[2])
	assert.Equal(t, "zach.mongohq.com", args[3])
	assert.Equal(t, "-u", args[4])
	assert.Equal(t, "user", args[5])
	assert.Equal(t, "-p", args[6])
	assert.Equal(t, "pass", args[7])
	assert.Equal(t, "-d", args[8])
	assert.Equal(t, "travisarchive", args[9])
	assert.Equal(t, "-c", args[10])
	assert.Equal(t, col, args[11])
	assert.Equal(t, "--out", args[12])
	assert.Equal(t, "file", args[13])
}
