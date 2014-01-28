package uploader

import (
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/bmizerany/assert"
)

func TestLocalUploader_Upload(t *testing.T) {
	destDir, _ := ioutil.TempDir("", "uploader")
	dest := filepath.Join(destDir, "test")

	f, err := ioutil.TempFile("", "uploader")
	assert.Equal(t, nil, err)

	err = ioutil.WriteFile(f.Name(), []byte("string"), 0644)
	assert.Equal(t, nil, err)

	uploader := Local{}
	err = uploader.Upload(dest, "application/json", f)
	assert.Equal(t, nil, err)

	c, err := ioutil.ReadFile(dest)
	assert.Equal(t, nil, err)
	assert.Equal(t, "string", string(c))
}
