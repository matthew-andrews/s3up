package etag

import (
	"crypto/md5"
	"fmt"
	"io"
	"os"
)

// With credit to http://dev.pawelsz.eu/2014/11/google-golang-compute-md5-of-file.html

func Compute(filePath string) (string, error) {
	var result []byte
	file, err := os.Open(filePath)
	if err != nil {
		return "", nil
	}
	defer file.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	result = hash.Sum([]byte{})
	return fmt.Sprintf("%x", result), nil
}
