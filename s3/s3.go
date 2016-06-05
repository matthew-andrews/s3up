package s3

import (
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/matthew-andrews/s3up/objects"
	"os"
)

func Upload(bucket string, files []objects.File) error {
	if len(files) < 1 {
		return errors.New("No files found for upload to S3.  (Directories are ignored)")
	}
	svc := s3.New(session.New())

	for _, file := range files {
		fmt.Printf("%s to %s\n", file.Location, file.Key)
		realFile, err := os.Open(file.Location)
		if err != nil {
			return errors.New(fmt.Sprintf("Could not open file: %s", file.Location))
		}
		defer realFile.Close()
		resp, err := svc.PutObject(&s3.PutObjectInput{
			Body:         realFile,
			Bucket:       &bucket,
			Key:          &file.Key,
			ContentType:  &file.ContentType,
			CacheControl: &file.CacheControl,
			ACL:          &file.ACL,
		})
		if err != nil {
			return errors.New(fmt.Sprintf("Failed to upload file to S3: %s", err.Error()))
		}
		fmt.Println(resp)
	}

	return nil
}
