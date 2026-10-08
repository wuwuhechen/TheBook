package utils

import (
	"testing"
)

func TestFindProjectRoot(t *testing.T) {
	root, err := FindProjectRoot()
	if err != nil {
		t.Fatalf("Failed to find project root: %v", err)
	}
	t.Logf("Project root: %s", root)
}
