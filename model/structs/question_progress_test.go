package structs

import "testing"

func TestNewQuestionProgress(t *testing.T) {
	progress := NewQuestionProgress(42)
	if progress.UserID != 42 || progress.CurrentQuestionID != 0 || progress.UpdatedAt.IsZero() {
		t.Fatalf("unexpected progress: %+v", progress)
	}
}
