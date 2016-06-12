package etag

import (
	"strings"
	"testing"
)

func TestEtagCompute(t *testing.T) {
	actual, err := Compute("../fixtures/one-file/my-file.txt")
	if err != nil {
		t.Fatalf("Unexpected error, %v\n", err)
	}
	expected := "f0ef7081e1539ac00ef5b761b4fb01b3"
	if actual != expected {
		t.Fatalf("Expected %s to equal %s\n", actual, expected)
	}
}

func TestEtagErrorsOnNonExistentFile(t *testing.T) {
	_, err := Compute("../fixtures/non-existent-file.txt")
	if err == nil {
		t.Fatalf("expected Compute to error")
	}
	if !strings.Contains(err.Error(), "no such file or directory") {
		t.Fatalf("expected %s to contain no such file or directory", err.Error())
	}
}
