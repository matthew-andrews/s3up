package s3client

import (
	"errors"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/matthew-andrews/s3up/objects"
	"testing"
)

type stubS3Service struct {
	UsedBucket string
}

func (stub stubS3Service) PutObject(input *s3.PutObjectInput) (*s3.PutObjectOutput, error) {
	var output *s3.PutObjectOutput

	if *input.Bucket != "my-fake-bucket" {
		return output, errors.New("Attempted to upload to the wrong bucket: " + *input.Bucket)
	}

	if *input.Key != "my-file.txt" {
		return output, errors.New("Attempted to upload to the wrong key: " + *input.Key)
	}

	return output, nil
}

func TestS3ClientUpload(t *testing.T) {
	stub := stubS3Service{}
	service := New(stub)
	err := service.Upload("my-fake-bucket", []objects.File{
		objects.File{
			Location:     "../fixtures/one-file/my-file.txt",
			Key:          "my-file.txt",
			Etag:         "123",
			ACL:          "public-read",
			CacheControl: "",
			ContentType:  "text/plain",
		},
	})
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}
}
