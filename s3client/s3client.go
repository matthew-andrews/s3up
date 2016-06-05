package s3client

import (
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/matthew-andrews/s3up/objects"
	"os"
)

type S3CompatibleInterface interface {
	PutObject(*s3.PutObjectInput) (*s3.PutObjectOutput, error)
}

type client struct {
	Service S3CompatibleInterface
}

func New(service S3CompatibleInterface) client {
	return client{
		Service: service,
	}
}

func (client *client) Upload(bucket string, files []objects.File) error {
	if len(files) < 1 {
		return errors.New("No files found for upload to S3.  (Directories are ignored)")
	}

	for _, file := range files {
		fmt.Printf("%s to %s\n", file.Location, file.Key)
		realFile, err := os.Open(file.Location)
		if err != nil {
			return errors.New(fmt.Sprintf("Could not open file: %s", file.Location))
		}
		defer realFile.Close()
		resp, err := client.Service.PutObject(&s3.PutObjectInput{
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
