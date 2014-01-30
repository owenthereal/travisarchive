package filestore

import (
	"fmt"
	"os"
	"time"
)

func New(adapter string) (FileStore, error) {
	u, ok := filestores[adapter]
	if !ok {
		return nil, fmt.Errorf("Unsupport adapter %s\n", adapter)
	}

	err := u.Init()
	if err != nil {
		return nil, err
	}

	return u, nil
}

var filestores = map[string]FileStore{
	"s3":    &S3{},
	"local": &Local{},
}

type File struct {
	Name string
	URI  string
	Time time.Time
}

func (f File) FormattedTime() string {
	return f.Time.Format("2006-01-02")
}

type FileStore interface {
	Init() error
	Upload(destPath, contentType string, f *os.File) error
	List(destPath string) ([]File, error)
}
