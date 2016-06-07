package objects

import (
	"reflect"
	"testing"
)

func TestGetFilesWithAFolderWithASingleFile(t *testing.T) {
	files, err := GetFiles([]string{
		"../fixtures/one-file/my-file.txt",
		// This should be skipped
		"../fixtures/one-file",
	}, 3, "prefix", "", "public-read")
	if err != nil {
		t.Fatalf("Unexpected error, %v\n", err)
	}

	if !reflect.DeepEqual(files, []File{
		File{
			ACL:          "public-read",
			CacheControl: "",
			ContentType:  "text/plain; charset=utf-8",
			Key:          "prefix/my-file.txt",
			Location:     "../fixtures/one-file/my-file.txt",
		},
	}) {
		t.Fatalf("%s did not match expected value", files)
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

	if !reflect.DeepEqual(files, []File{
		File{
			ACL:          "public-read",
			CacheControl: "",
			ContentType:  "text/plain; charset=utf-8",
			Key:          "subsubfolder/bottom-file.txt",
			Location:     "../fixtures/subfolders/subsubfolder/bottom-file.txt",
		},
		File{
			ACL:          "public-read",
			CacheControl: "",
			ContentType:  "text/plain; charset=utf-8",
			Key:          "top-file.txt",
			Location:     "../fixtures/subfolders/top-file.txt",
		},
	}) {
		t.Fatalf("%s did not match expected value", files)
	}

}
