package ls

import (
	"errors"
	"fmt"
	"testing"
)

func TestGetAll(t *testing.T) {
	files, err := GetFiles("../fixtures/one-file")
	if err != nil {
		t.Fatalf("Unexpected error, %v\n", err)
	}

	if err = fileSlicesAreEquivalent(files, []File{
		File{
			Filename:    "my-file.txt",
			ContentType: "text/plain; charset=utf-8",
			Etag:        "TODO",
		},
	}); err != nil {
		t.Fatalf("Unexpected error, %v\n", err)
	}
}

func fileSlicesAreEquivalent(expected []File, actual []File) error {
	if len(expected) != len(actual) {
		return errors.New(fmt.Sprintf("returned slice should be length %d", len(expected)))
	}

	for i, _ := range expected {
		if expected[i] != actual[i] {
			return errors.New(fmt.Sprintf("item at index %d should match %s but was %s", i, expected[i], actual[i]))
		}
	}

	return nil
}
