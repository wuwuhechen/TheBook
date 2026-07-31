package utils_test

import (
	"TheBook/utils"
	"testing"
)

func TestFindProjectRoot(t *testing.T) {
	root, err := utils.FindProjectRoot()
	if err != nil {
		t.Fatalf("Failed to find project root: %v", err)
	}
	t.Logf("Project root: %s", root)
}
