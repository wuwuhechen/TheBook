package bank

import "testing"

// TestQuestionProgressBankFindAndUpsert 验证内存进度查询和更新。
func TestQuestionProgressBankFindAndUpsert(t *testing.T) {
	bank := NewQuestionProgressBank()
	newProgress, err := bank.FindByUserID(1)
	if err != nil || newProgress.UserID != 1 || newProgress.CurrentQuestionID != 0 {
		t.Fatalf("unexpected new progress: %+v, %v", newProgress, err)
	}
	if err := bank.Upsert(&QuestionProgress{UserID: 1, CurrentQuestionID: 12}); err != nil {
		t.Fatalf("upsert: %v", err)
	}
	stored, err := bank.FindByUserID(1)
	if err != nil || stored.CurrentQuestionID != 12 {
		t.Fatalf("unexpected stored progress: %+v, %v", stored, err)
	}
	if err := bank.Upsert(&QuestionProgress{UserID: 1, CurrentQuestionID: 23}); err != nil {
		t.Fatalf("second upsert: %v", err)
	}
	updated, _ := bank.FindByUserID(1)
	if updated.CurrentQuestionID != 23 {
		t.Fatalf("expected progress replacement, got %+v", updated)
	}
}
