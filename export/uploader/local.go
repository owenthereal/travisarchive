package uploader

import (
	"io"
	"os"
)

type Local struct {
}

func (LocalUploader *Local) Init() error {
	return nil
}

func (LocalUploader *Local) Upload(destPath, contentType string, f *os.File) error {
	destFile, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = io.Copy(destFile, f)

	return err
}
