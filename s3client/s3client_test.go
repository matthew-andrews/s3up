package s3client

import (
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/matthew-andrews/s3up/objects"
	"testing"
)

// Stubbing

var lastPutObjectInput *s3.PutObjectInput

type stubS3Service struct{}

func (stub stubS3Service) PutObject(input *s3.PutObjectInput) (*s3.PutObjectOutput, error) {
	lastPutObjectInput = input
	var putObjectOutput *s3.PutObjectOutput
	return putObjectOutput, nil
}
func reset() {
	lastPutObjectInput = nil
}

// Sample data

func sampleFiles() []objects.File {
	return []objects.File{
		objects.File{
			Location:     "../fixtures/one-file/my-file.txt",
			Key:          "my-file.txt",
			Etag:         "123",
			ACL:          "public-read",
			CacheControl: "",
			ContentType:  "text/plain",
		},
	}
}

func TestS3ClientUpload(t *testing.T) {
	reset()
	stub := stubS3Service{}
	service := New(stub)
	err := service.Upload("my-fake-bucket", sampleFiles())
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
