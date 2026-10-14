package practice_test

import (
	"TheBook/model"
	"TheBook/test/gin_handler_test/shared"
	"net/http"
	"strings"
	"testing"
)

// TestPracticePostInitAnswerAndSubmit 验证套题初始化、保存答案和提交 Handler。
func TestPracticePostInitAnswerAndSubmit(t *testing.T) {
	f := shared.NewFixture(t)
	w := shared.Serve(f, f.Request(http.MethodPost, "/practice/init", []byte(`{`)))
	if w.Code != http.StatusBadRequest {
		t.Fatalf("invalid practice init: got %d", w.Code)
	}
	w = shared.Serve(f, f.Request(http.MethodPost, "/practice/init", shared.JSONBody(t, model.Request{PracticeSize: 2})))
	if w.Code != http.StatusFound || !strings.HasPrefix(w.Header().Get("Location"), "/practice/") {
		t.Fatalf("practice init failed: %d", w.Code)
	}

	p := shared.NewPractice(t, f, 501, 1, 1, 2)
	payload := shared.JSONBody(t, model.Request{PracticeID: p.ID, QuestionID: 1, Choice: 0})
	w = shared.Serve(f, f.RequestAs(t, 2, http.MethodPost, "/practice/answer", payload))
	if w.Code != http.StatusNotFound {
		t.Fatalf("non-owner answer: got %d", w.Code)
	}
	w = shared.Serve(f, f.Request(http.MethodPost, "/practice/answer", payload))
	if w.Code != http.StatusOK || p.Answers[1] != 0 {
		t.Fatalf("first answer save failed: %+v", p.Answers)
	}
	payload = shared.JSONBody(t, model.Request{PracticeID: p.ID, QuestionID: 1, Choice: 1})
	w = shared.Serve(f, f.Request(http.MethodPost, "/practice/answer", payload))
	if w.Code != http.StatusOK || len(p.Answers) != 1 || p.Answers[1] != 1 {
		t.Fatalf("answer was not overwritten: %+v", p.Answers)
	}

	w = shared.Serve(f, f.Request(http.MethodPost, "/practice/501/submit", nil))
	if w.Code != http.StatusOK || !strings.Contains(w.Body.String(), `"correct_count":1`) {
		t.Fatalf("practice submit failed: %d: %s", w.Code, w.Body.String())
	}
	if !p.Completed {
		t.Fatal("submitted practice not marked complete")
	}
	w = shared.Serve(f, f.Request(http.MethodPost, "/practice/501/submit", nil))
	if w.Code != http.StatusConflict {
		t.Fatalf("duplicate submit: got %d", w.Code)
	}
	wrong, err := f.Server.WQ.ListWrongQuestionsByUser(1, 10, 0)
	if err != nil || len(wrong) != 1 || wrong[0].QuestionID != 2 {
		t.Fatalf("unanswered question missing from wrong bank: %+v, %v", wrong, err)
	}
}

// TestPracticeSubmitPersistsOwnedResult 验证只有套题所有者可提交，且提交结果会被隔离持久化。
func TestPracticeSubmitPersistsOwnedResult(t *testing.T) {
	f := shared.NewFixture(t)
	p := shared.NewPractice(t, f, 502, 1, 1, 2)
	p.Answers[1] = 1

	w := shared.Serve(f, f.RequestAs(t, 2, http.MethodPost, "/practice/502/submit", nil))
	if w.Code != http.StatusNotFound {
		t.Fatalf("non-owner submit: got %d: %s", w.Code, w.Body.String())
	}
	if p.Completed {
		t.Fatal("non-owner submit marked practice as completed")
	}
	if _, err := f.Server.PM.GetRecordByPracticeID(1, p.ID); err == nil {
		t.Fatal("non-owner submit unexpectedly persisted a practice record")
	}

	w = shared.Serve(f, f.Request(http.MethodPost, "/practice/502/submit", nil))
	if w.Code != http.StatusOK {
		t.Fatalf("owner submit failed: got %d: %s", w.Code, w.Body.String())
	}
	record, err := f.Server.PM.GetRecordByPracticeID(1, p.ID)
	if err != nil {
		t.Fatalf("load persisted record: %v", err)
	}
	if record.TotalQuestions != 2 || record.CorrectCount != 1 || record.WrongCount != 1 || len(record.Answers) != 2 {
		t.Fatalf("unexpected persisted record: %+v", record)
	}
	if _, err := f.Server.PM.GetRecordByPracticeID(2, p.ID); err == nil {
		t.Fatal("another user can read the owner's practice record")
	}

	ownerWrongQuestions, err := f.Server.WQ.ListWrongQuestionsByUser(1, 10, 0)
	if err != nil || len(ownerWrongQuestions) != 1 || ownerWrongQuestions[0].QuestionID != 2 {
		t.Fatalf("owner wrong questions: %+v, %v", ownerWrongQuestions, err)
	}
	otherWrongQuestions, err := f.Server.WQ.ListWrongQuestionsByUser(2, 10, 0)
	if err != nil || len(otherWrongQuestions) != 0 {
		t.Fatalf("wrong questions leaked to another user: %+v, %v", otherWrongQuestions, err)
	}
}
