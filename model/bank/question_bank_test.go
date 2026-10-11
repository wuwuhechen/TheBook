package bank

import "testing"

// TestQuestionBankAddReplaceAndRead 验证内存题库的添加、替换和读取。
func TestQuestionBankAddReplaceAndRead(t *testing.T) {
	bank := NewQuestionBank()
	bank.AddQuestion(Question{ID: 2, Question: "old"})
	bank.AddQuestion(Question{ID: 1, Question: "first"})
	bank.AddQuestion(Question{ID: 2, Question: "replacement"})

	if bank.GetTotalCount() != 2 {
		t.Fatalf("expected 2 unique questions, got %d", bank.GetTotalCount())
	}
	question, err := bank.GetQuestion(2)
	if err != nil || question.Question != "replacement" {
		t.Fatalf("unexpected replacement result: %+v, %v", question, err)
	}
	if _, err := bank.GetQuestion(999); err == nil {
		t.Fatal("expected missing question error")
	}

	sorted := bank.GetAllQuestionsSorted()
	if len(sorted) != 2 || sorted[0].ID != 1 || sorted[1].ID != 2 {
		t.Fatalf("unexpected sorted questions: %+v", sorted)
	}
	ids := bank.GetALLQuestionIDs()
	ids[0] = 999
	if bank.GetALLQuestionIDs()[0] == 999 {
		t.Fatal("returned IDs must be a copy")
	}
}

// TestQuestionBankRandomQuestion 验证内存题库的随机取题。
func TestQuestionBankRandomQuestion(t *testing.T) {
	empty := NewQuestionBank()
	if _, err := empty.GetRandomQuestionID(); err == nil {
		t.Fatal("expected error for empty bank")
	}

	bank := NewQuestionBank()
	bank.AddQuestion(Question{ID: 10})
	bank.AddQuestion(Question{ID: 20})
	for range 20 {
		question, err := bank.GetRandomQuestionID()
		if err != nil || (question.ID != 10 && question.ID != 20) {
			t.Fatalf("unexpected random question: %+v, %v", question, err)
		}
	}
}
