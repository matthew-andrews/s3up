package objects

import (
	"testing"
)

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
