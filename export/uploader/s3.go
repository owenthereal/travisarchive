package uploader

import (
	"fmt"
	"os"

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

func (s3Uploader *S3) Init() error {
	s3Region := os.Getenv("S3_REGION")
	region, ok := aws.Regions[s3Region]
	if !ok {
		return fmt.Errorf("Fail to find S3 region %s\n", s3Region)
	}

	auth := aws.Auth{AccessKey: os.Getenv("AWS_ACCESS_KEY_ID"), SecretKey: os.Getenv("AWS_SECRET_KEY")}
	s := s3.New(auth, region)
	s3Uploader.Bucket = s.Bucket(os.Getenv("S3_BUCKET"))
	s3Uploader.BufferSize = defaultS3BufferSize

	return nil
}

func (s3Uploader *S3) Upload(destPath, contentType string, f *os.File) error {
	writer, err := s3Uploader.Bucket.InitMulti(destPath, contentType, s3.PublicRead)
	if err != nil {
		return err
	}

	parts, err := writer.PutAll(f, s3Uploader.BufferSize)
	if err != nil {
		return err
	}

	return writer.Complete(parts)
}
