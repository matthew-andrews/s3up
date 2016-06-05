package s3client

import (
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/matthew-andrews/s3up/objects"
	"os"
)

type s3Interface interface {
	PutObject(*s3.PutObjectInput) (*s3.PutObjectOutput, error)
	HeadObject(*s3.HeadObjectInput) (*s3.HeadObjectOutput, error)
}

type client struct {
	Service s3Interface
}

func New(service s3Interface) client {
	return client{
		Service: service,
	}
}

func (client client) Upload(bucket string, files []objects.File) error {
	if len(files) < 1 {
		return errors.New("No files found for upload to S3.  (Directories are ignored)")
	}

	for _, file := range files {
		if err := client.upload(bucket, file); err != nil {
			return err
		}
	}

	return nil
}

func (client client) upload(bucket string, file objects.File) error {
	// HeadObject
	headResp, err := client.Service.HeadObject(&s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(file.Key),
	})

	if err != nil {
		return errors.New(fmt.Sprintf("Head request to S3 object failed: %s", err))
	}

	if aws.StringValue(headResp.ETag) == "\""+file.ETag+"\"" {
		fmt.Printf("Unchanged, skipping: %s\n", file.Key)
	} else {
		// PutObject
		realFile, err := os.Open(file.Location)
		if err != nil {
			return errors.New(fmt.Sprintf("Could not open file: %s", file.Location))
		}
		defer realFile.Close()
		_, err = client.Service.PutObject(&s3.PutObjectInput{
			Body:         realFile,
			Bucket:       aws.String(bucket),
			Key:          aws.String(file.Key),
			ContentType:  aws.String(file.ContentType),
			CacheControl: aws.String(file.CacheControl),
			ACL:          aws.String(file.ACL),
		})
		if err != nil {
			return errors.New(fmt.Sprintf("Failed to upload file to S3: %s", err.Error()))
		}
		fmt.Printf("Successfully uploaded: %s", file.Key)
	}
	return nil
}
