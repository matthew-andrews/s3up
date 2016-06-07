package uploader

import (
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/matthew-andrews/s3up/objects"
	"github.com/matthew-andrews/s3up/s3client"
	"sync"
)

func Upload(bucket string, files []objects.File, concurrency int) error {
	if len(files) < 1 {
		return errors.New("No files found for upload to S3.  (Directories are ignored)")
	}

	var sem = make(chan bool, concurrency)
	var wg sync.WaitGroup
	service := s3client.New(s3.New(session.New()))

	for _, file := range files {
		wg.Add(1)
		sem <- true
		go func(file objects.File) {
			defer wg.Done()
			defer func() { <-sem }()
			if err := service.UploadFile(bucket, file); err != nil {
				fmt.Println("Failed to upload", file, err)
				// TODO handle error
			}
		}(file)
	}

	wg.Wait()
	return nil
}
