package s3client

import (
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/matthew-andrews/s3up/objects"
	"testing"
)

type stubS3Service struct{}

var lastPutObjectInput *s3.PutObjectInput

func (stub stubS3Service) PutObject(input *s3.PutObjectInput) (*s3.PutObjectOutput, error) {
	lastPutObjectInput = input
	var putObjectOutput *s3.PutObjectOutput
	return putObjectOutput, nil
}

func reset() {
	lastPutObjectInput = nil
}

func TestS3ClientUpload(t *testing.T) {
	reset()
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

	if *lastPutObjectInput.Bucket != "my-fake-bucket" {
		t.Fatalf("Attempted to upload to the wrong bucket: %s", *lastPutObjectInput.Bucket)
	}

	if *lastPutObjectInput.Key != "my-file.txt" {
		t.Fatalf("Attempted to upload to the wrong key: %s", *lastPutObjectInput.Key)
	}
}
