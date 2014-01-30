package filestore

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

type Local struct {
}

func (l *Local) Init() error {
	return nil
}

func (l *Local) Upload(destPath, contentType string, f *os.File) error {
	destFile, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = io.Copy(destFile, f)

	return err
}

func (l *Local) List(destPath string) (files []File, err error) {
	info, err := ioutil.ReadDir(destPath)
	if err != nil {
		return
	}

	for _, i := range info {
		file := File{Name: i.Name(), Time: i.ModTime(), URI: filepath.Join(destPath, i.Name())}
		files = append(files, file)
	}

	return
}
