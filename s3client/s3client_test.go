package s3client

import (
	"errors"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/matthew-andrews/s3up/objects"
	"testing"
)

// Stubbing

var lastUploadInput *s3manager.UploadInput
var lastHeadObjectInput *s3.HeadObjectInput
var lastCopyObjectInput *s3.CopyObjectInput

type stubS3Uploader struct{}

func (stub stubS3Uploader) Upload(input *s3manager.UploadInput, options ...func(*s3manager.Uploader)) (*s3manager.UploadOutput, error) {
	lastUploadInput = input
	var uploadOutput *s3manager.UploadOutput
	return uploadOutput, nil
}

type stubS3Service struct{}

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

	if aws.StringValue(input.Key) == "my-new-file.txt" {
		return &s3.HeadObjectOutput{}, awserr.New("", "status code: 404, request id: AAAAAAAAAAAAAAAA", errors.New("status code: 404, request id: AAAAAAAAAAAAAAAA"))
	}

	return &s3.HeadObjectOutput{}, nil
}
func testUpload(file objects.File) error {
	lastUploadInput = nil
	lastHeadObjectInput = nil
	lastCopyObjectInput = nil
	service := New(stubS3Service{}, stubS3Uploader{})
	return service.UploadFile("my-fake-bucket", file)
}

// Sample data

func TestS3ClientUploadFile(t *testing.T) {
	err := testUpload(objects.File{
		Location:     "../fixtures/one-file/my-file.txt",
		Key:          "my-new-file.txt",
		ACL:          "public-read",
		CacheControl: "",
		ContentType:  "text/plain",
	})
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}

	if aws.StringValue(lastUploadInput.Bucket) != "my-fake-bucket" {
		t.Fatalf("Attempted to upload to the wrong bucket: %s", aws.StringValue(lastUploadInput.Bucket))
	}

	if aws.StringValue(lastUploadInput.Key) != "my-new-file.txt" {
		t.Fatalf("Attempted to upload to the wrong key: %s", aws.StringValue(lastUploadInput.Key))
	}
}

func TestHeadsBeforePuts(t *testing.T) {
	err := testUpload(objects.File{
		Location:     "../fixtures/one-file/my-file.txt",
		Key:          "my-file.txt",
		ACL:          "public-read",
		CacheControl: "",
		ContentType:  "text/plain",
	})
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}

	if lastHeadObjectInput == nil {
		t.Fatalf("S3Client should make a HeadObject request to the S3 object before deciding to upload")
	}
	if lastUploadInput != nil {
		t.Fatalf("S3Client should make not have made a Upload request to the S3 object if the file hasn't changed")
	}
}

func TestUpdatesMetadataIfThatIsAllThatHasChanged(t *testing.T) {
	err := testUpload(objects.File{
		Location:     "../fixtures/one-file/my-file.txt",
		Key:          "my-file-with-different-metadata.txt",
		ACL:          "public-read",
		CacheControl: "",
		ContentType:  "text/plain",
	})
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}

	if lastHeadObjectInput == nil {
		t.Fatalf("S3Client should make a HeadObject request to the S3 object before deciding to upload")
	}
	if aws.StringValue(lastCopyObjectInput.ContentType) != "text/plain" {
		t.Fatalf("S3Client should have CopyObject request to the S3 object to update the metadata if it has changed")
	}
	if lastUploadInput != nil {
		t.Fatalf("S3Client should make not have made a Upload request with a Body to the S3 object if the file hasn't changed")
	}
	if lastCopyObjectInput == nil {
		t.Fatalf("S3Client should not have made a CopyObject request if only the file's metadata has changed")
	}
}
