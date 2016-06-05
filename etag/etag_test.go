package etag

import (
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
