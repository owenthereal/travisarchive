package filestore

import (
	"fmt"
	"github.com/jingweno/travisarchive/util"
	"os"
	"path/filepath"
	"strings"

	"launchpad.net/goamz/aws"
	"launchpad.net/goamz/s3"
)

const (
	defaultS3BufferSize = 5 * 1024 * 1024
)

type S3 struct {
	Bucket     *s3.Bucket
	BufferSize int64
}

func (s3FS *S3) Init() error {
	s3Region := os.Getenv("S3_REGION")
	region, ok := aws.Regions[s3Region]
	if !ok {
		return fmt.Errorf("Fail to find S3 region %s\n", s3Region)
	}

	auth := aws.Auth{AccessKey: os.Getenv("AWS_ACCESS_KEY_ID"), SecretKey: os.Getenv("AWS_SECRET_KEY")}
	s := s3.New(auth, region)
	s3FS.Bucket = s.Bucket(os.Getenv("S3_BUCKET"))
	s3FS.BufferSize = defaultS3BufferSize

	return nil
}

func (s3FS *S3) Upload(destPath, contentType string, f *os.File) error {
	writer, err := s3FS.Bucket.InitMulti(destPath, contentType, s3.PublicRead)
	if err != nil {
		return err
	}

	parts, err := writer.PutAll(f, s3FS.BufferSize)
	if err != nil {
		return err
	}

	return writer.Complete(parts)
}

func (s3FS *S3) List(destPath string) (files []File, err error) {
	destPath = strings.TrimPrefix(destPath, "/")
	destPath = destPath + "/"

	results, err := s3FS.Bucket.List(destPath, "/", "", 1000)
	for _, c := range results.Contents {
		name := strings.TrimPrefix(c.Key, "builds/")
		if name == "" {
			continue
		}

		nameToParse := strings.TrimSuffix(name, filepath.Ext(name))
		t, err := util.ParseBuildTime(nameToParse)
		if err != nil {
			continue
		}

		uri := fmt.Sprintf("https://s3.amazonaws.com/travisarchive/%s", c.Key)
		file := File{Name: name, Time: t, URI: uri}
		files = append(files, file)
	}

	return
}
