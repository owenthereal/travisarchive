package main

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type Archiver struct {
	File string
}

func (a *Archiver) Archive() (out string, err error) {
	out = fmt.Sprintf("%s.zip", strings.TrimSuffix(a.File, filepath.Ext(a.File)))
	file, err := os.OpenFile(out, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return
	}
	defer file.Close()

	zw := zip.NewWriter(file)
	defer zw.Close()

	err = a.addFile(a.File, zw)

	return
}

func (a *Archiver) addFile(filename string, zw *zip.Writer) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("failed opening %s: %s", filename, err)
	}
	defer file.Close()

	wr, err := zw.Create(filepath.Base(filename))
	if err != nil {
		return fmt.Errorf("failed creating entry for %s in zip file: %s", filename, err)
	}

	if _, err := io.Copy(wr, file); err != nil {
		return fmt.Errorf("failed writing %s to zip: %s", filename, err)
	}

	return nil
}
