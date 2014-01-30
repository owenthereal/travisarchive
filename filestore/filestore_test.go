package filestore

import (
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/bmizerany/assert"
)

func TestLocal_Upload(t *testing.T) {
	destDir, _ := ioutil.TempDir("", "localfs")
	dest := filepath.Join(destDir, "test")

	f, err := ioutil.TempFile("", "localfs")
	assert.Equal(t, nil, err)

	err = ioutil.WriteFile(f.Name(), []byte("string"), 0644)
	assert.Equal(t, nil, err)

	fs := Local{}
	err = fs.Upload(dest, "application/json", f)
	assert.Equal(t, nil, err)

	c, err := ioutil.ReadFile(dest)
	assert.Equal(t, nil, err)
	assert.Equal(t, "string", string(c))
}
