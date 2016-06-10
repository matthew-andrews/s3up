package uploader

import (
	"errors"
	"github.com/matthew-andrews/s3up/objects"
	"strings"
	"testing"
	"time"
)

type stubS3Client struct{}

func (stub stubS3Client) UploadFile(bucket string, object objects.File) error {
	if object.Key == "error" {
		return errors.New("synthetic error")
	}
	time.Sleep(50 * time.Millisecond)
	return nil
}

func uploadThreeFilesWithConcurrency(concurrency int) int64 {
	startTime := time.Now()
	Upload(stubS3Client{}, "", make([]objects.File, 3), concurrency)
	duration := time.Since(startTime).Nanoseconds()
	return int64(duration / int64(time.Millisecond))
}

func TestOneAtATime(t *testing.T) {
	duration := uploadThreeFilesWithConcurrency(1)
	if duration < 100 {
		t.Fatalf("uploader was too quick.  3 times 50ms one at a time can't be less than 100ms.  but it was %v", duration)
	}
}

func TestThreeAtATime(t *testing.T) {
	duration := uploadThreeFilesWithConcurrency(3)
	if duration > 100 {
		t.Fatalf("uploader was too slow.  3 times 50ms three at a time can't be more than 100ms.  but it was %v", duration)
	}
}

func TestErrorsAreEscalated(t *testing.T) {
	errors := Upload(stubS3Client{}, "", []objects.File{
		objects.File{Key: "error"},
		objects.File{},
		objects.File{Key: "error"},
	}, 2)
	if len(errors) != 2 {
		t.Fatal("Uploader should have returned 2 errors")
	}
}

func TestNoFiles(t *testing.T) {
	errs := Upload(stubS3Client{}, "", make([]objects.File, 0), 1)
	if strings.Contains(errs[0].Error(), "No files found") == false {
		t.Fatal("The error that was expected was not thrown")
	}
}
