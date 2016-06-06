package uploader

import (
	"errors"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/matthew-andrews/s3up/objects"
	"github.com/matthew-andrews/s3up/s3client"
)

func Upload(bucket string, files []objects.File) error {
	service := s3client.New(s3.New(session.New()))

	if len(files) < 1 {
		return errors.New("No files found for upload to S3.  (Directories are ignored)")
	}

	for _, file := range files {
		if err := service.UploadFile(bucket, file); err != nil {
			return err
		}
	}

	return nil
}
