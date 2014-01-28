package uploader

import (
	"fmt"
	"os"
)

func New(adapter string) (Uploader, error) {
	u, ok := uploaders[adapter]
	if !ok {
		return nil, fmt.Errorf("Unsupport adapter %s\n", adapter)
	}

	err := u.Init()
	if err != nil {
		return nil, err
	}

	return u, nil
}

var uploaders = map[string]Uploader{
	"s3":    &S3{},
	"local": &Local{},
}

type Uploader interface {
	Init() error
	Upload(destPath, contentType string, f *os.File) error
}
