package uploader

import (
	"errors"
	"github.com/matthew-andrews/s3up/objects"
	"sync"
)

type s3ClientInterface interface {
	UploadFile(string, objects.File) error
}

func Upload(service s3ClientInterface, bucket string, files []objects.File, concurrency int) []error {
	ec := make(chan error, len(files))
	if len(files) < 1 {
		return []error{errors.New("No files found for upload to S3.  (Directories are ignored)")}
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
				ec <- err
			}
		}(file)
	}

	wg.Wait()
	close(ec)
	var errs []error
	if len(ec) > 0 {
		for err := range ec {
			errs = append(errs, err)
		}
	}

	return errs
}
