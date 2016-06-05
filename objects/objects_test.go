package objects

import (
	"errors"
	"fmt"
	"testing"
)

func TestGetFilesWithAFolderWithASingleFile(t *testing.T) {
	files, err := GetFiles([]string{"../fixtures/one-file/my-file.txt"}, 3, "prefix", "", "public-read")
	if err != nil {
		t.Fatalf("Unexpected error, %v\n", err)
	}

	if err = fileSlicesAreEquivalent(files, []File{
		File{
			ACL:          "public-read",
			CacheControl: "",
			ContentType:  "text/plain; charset=utf-8",
			Etag:         "TODO",
			Key:          "prefix/my-file.txt",
			Location:     "../fixtures/one-file/my-file.txt",
		},
	}); err != nil {
		t.Fatalf("Unexpected error, %v\n", err)
	}
}

func TestGetFilesWithAFolderWithSubfolders(t *testing.T) {
	files, err := GetFiles([]string{
		"../fixtures/subfolders/subsubfolder/bottom-file.txt",
		"../fixtures/subfolders/top-file.txt",
	}, 3, "", "", "public-read")
	if err != nil {
		t.Fatalf("Unexpected error, %v\n", err)
	}

	if err = fileSlicesAreEquivalent(files, []File{
		File{
			ACL:          "public-read",
			CacheControl: "",
			ContentType:  "text/plain; charset=utf-8",
			Etag:         "TODO",
			Key:          "subsubfolder/bottom-file.txt",
			Location:     "../fixtures/subfolders/subsubfolder/bottom-file.txt",
		},
		File{
			ACL:          "public-read",
			CacheControl: "",
			ContentType:  "text/plain; charset=utf-8",
			Etag:         "TODO",
			Key:          "top-file.txt",
			Location:     "../fixtures/subfolders/top-file.txt",
		},
	}); err != nil {
		t.Fatalf("Unexpected error, %v\n", err)
	}

}

func TestStripFromName(t *testing.T) {
	stripMap := map[int]string{
		0: "../fixtures/one-file/my-file.txt",
		1: "fixtures/one-file/my-file.txt",
		2: "one-file/my-file.txt",
		3: "my-file.txt",
	}
	for strip, expected := range stripMap {
		if actual := StripFromName("../fixtures/one-file/my-file.txt", strip); actual != expected {
			t.Fatalf("expected %s to equal %s", actual, expected)
		}
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
