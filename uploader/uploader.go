package uploader

import (
	"errors"
	"fmt"
	"github.com/matthew-andrews/s3up/objects"
	"sync"
)

type s3ClientInterface interface {
	UploadFile(string, objects.File) error
}

func Upload(service s3ClientInterface, bucket string, files []objects.File, concurrency int) error {
	if len(files) < 1 {
		return errors.New("No files found for upload to S3.  (Directories are ignored)")
	}

	var sem = make(chan bool, concurrency)
	var wg sync.WaitGroup

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
