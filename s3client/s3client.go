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
		if err := client.uploadFile(bucket, file); err != nil {
			return err
		}
	}

	return nil
}

func (client client) uploadFile(bucket string, file objects.File) error {
	// HeadObject
	headResp, err := client.Service.HeadObject(&s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(file.Key),
	})

	if err != nil {
		return errors.New(fmt.Sprintf("Head request to S3 object failed: %s", err))
	}

	putObjectInput := &s3.PutObjectInput{
		Bucket:       aws.String(bucket),
		Key:          aws.String(file.Key),
		ContentType:  aws.String(file.ContentType),
		CacheControl: aws.String(file.CacheControl),
		ACL:          aws.String(file.ACL),
	}

	// Add Body, if ETags mismatch
	if aws.StringValue(headResp.ETag) != "\""+file.ETag+"\"" {
		realFile, err := os.Open(file.Location)
		if err != nil {
			return errors.New(fmt.Sprintf("Could not open file: %s", file.Location))
		}
		defer realFile.Close()
		putObjectInput.Body = realFile
	}

	// If the file has changed in some way, update it
	if putObjectInput.Body != nil || aws.StringValue(headResp.CacheControl) != file.CacheControl || aws.StringValue(headResp.ContentType) != file.ContentType {
		_, err := client.Service.PutObject(putObjectInput)
		if err != nil {
			return errors.New(fmt.Sprintf("Failed to update file: %s", err.Error()))
		}
		if putObjectInput.Body == nil {
			fmt.Printf("Successfully updated just metadata of: %s\n", file.Key)
		} else {
			fmt.Printf("Successfully uploaded: %s\n", file.Key)
		}
	} else {
		fmt.Printf("Unchanged, skipping: %s\n", file.Key)
	}

	return nil
}
