package utils

import (
	"testing"
)

// TestFindProjectRoot 验证能够定位项目根目录。
func TestFindProjectRoot(t *testing.T) {
	root, err := FindProjectRoot()
	if err != nil {
		t.Fatalf("Failed to find project root: %v", err)
	}
	t.Logf("Project root: %s", root)
}
