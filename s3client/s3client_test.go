package s3client

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/matthew-andrews/s3up/objects"
	"testing"
)

// Stubbing

var lastPutObjectInput *s3.PutObjectInput
var lastHeadObjectInput *s3.HeadObjectInput
var lastCopyObjectInput *s3.CopyObjectInput

type stubS3Service struct{}

func (stub stubS3Service) PutObject(input *s3.PutObjectInput) (*s3.PutObjectOutput, error) {
	lastPutObjectInput = input
	var putObjectOutput *s3.PutObjectOutput
	return putObjectOutput, nil
}
func (stub stubS3Service) CopyObject(input *s3.CopyObjectInput) (*s3.CopyObjectOutput, error) {
	lastCopyObjectInput = input
	var copyObjectOutput *s3.CopyObjectOutput
	return copyObjectOutput, nil
}
func (stub stubS3Service) HeadObject(input *s3.HeadObjectInput) (*s3.HeadObjectOutput, error) {
	lastHeadObjectInput = input

	if aws.StringValue(input.Key) == "my-file.txt" {
		return &s3.HeadObjectOutput{
			CacheControl: aws.String(""),
			ContentType:  aws.String("text/plain"),
			ETag:         aws.String("\"f0ef7081e1539ac00ef5b761b4fb01b3\""),
		}, nil
	}

	if aws.StringValue(input.Key) == "my-file-with-different-metadata.txt" {
		return &s3.HeadObjectOutput{
			CacheControl: aws.String(""),
			ContentType:  aws.String("text/html"),
			ETag:         aws.String("\"f0ef7081e1539ac00ef5b761b4fb01b3\""),
		}, nil
	}

	return &s3.HeadObjectOutput{}, nil
}
func reset() {
	lastPutObjectInput = nil
	lastHeadObjectInput = nil
}

// Sample data

func TestS3ClientUpload(t *testing.T) {
	reset()
	stub := stubS3Service{}
	service := New(stub)
	err := service.Upload("my-fake-bucket", []objects.File{
		objects.File{
			Location:     "../fixtures/one-file/my-file.txt",
			Key:          "my-new-file.txt",
			ETag:         "f0ef7081e1539ac00ef5b761b4fb01b3",
			ACL:          "public-read",
			CacheControl: "",
			ContentType:  "text/plain",
		},
	})
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}

	if aws.StringValue(lastPutObjectInput.Bucket) != "my-fake-bucket" {
		t.Fatalf("Attempted to upload to the wrong bucket: %s", aws.StringValue(lastPutObjectInput.Bucket))
	}

	if aws.StringValue(lastPutObjectInput.Key) != "my-new-file.txt" {
		t.Fatalf("Attempted to upload to the wrong key: %s", aws.StringValue(lastPutObjectInput.Key))
	}
}

func TestHeadsBeforePuts(t *testing.T) {
	reset()
	stub := stubS3Service{}
	service := New(stub)
	err := service.Upload("my-fake-bucket", []objects.File{
		objects.File{
			Location:     "../fixtures/one-file/my-file.txt",
			Key:          "my-file.txt",
			ETag:         "f0ef7081e1539ac00ef5b761b4fb01b3",
			ACL:          "public-read",
			CacheControl: "",
			ContentType:  "text/plain",
		},
	})
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}

	if lastHeadObjectInput == nil {
		t.Fatalf("S3Client should make a HeadObject request to the S3 object before deciding to upload")
	}
	if lastPutObjectInput != nil {
		t.Fatalf("S3Client should make not have made a PutObject request to the S3 object if the file hasn't changed")
	}
}

func TestUpdatesMetadataIfThatIsAllThatHasChanged(t *testing.T) {
	reset()
	stub := stubS3Service{}
	service := New(stub)
	err := service.Upload("my-fake-bucket", []objects.File{
		objects.File{
			Location:     "../fixtures/one-file/my-file.txt",
			Key:          "my-file-with-different-metadata.txt",
			ETag:         "f0ef7081e1539ac00ef5b761b4fb01b3",
			ACL:          "public-read",
			CacheControl: "",
			ContentType:  "text/plain",
		},
	})
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}

	if lastHeadObjectInput == nil {
		t.Fatalf("S3Client should make a HeadObject request to the S3 object before deciding to upload")
	}
	if aws.StringValue(lastCopyObjectInput.ContentType) != "text/plain" {
		t.Fatalf("S3Client should have PutObject request to the S3 object to update the metadata if it has changed")
	}
	if lastPutObjectInput != nil {
		t.Fatalf("S3Client should make not have made a PutObject request with a Body to the S3 object if the file hasn't changed")
	}
	if lastCopyObjectInput == nil {
		t.Fatalf("S3Client should not have made a CopyObject request if only the file's metadata has changed")
	}
}
