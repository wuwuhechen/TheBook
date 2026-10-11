package practice_test

import (
	"TheBook/model"
	"TheBook/test/gin_handler_test/shared"
	"net/http"
	"strings"
	"testing"
	"time"
)

// TestPracticeGetNavigationAndHistory 验证套题导航和历史记录 Handler。
func TestPracticeGetNavigationAndHistory(t *testing.T) {
	f := shared.NewFixture(t)
	p := shared.NewPractice(t, f, 101, 1, 1, 2)
	p.Duration = 5 * time.Minute
	w := shared.Serve(f, f.Request(http.MethodGet, "/practice/101?direction=next", nil))
	if w.Code != http.StatusOK || p.CurrentIndex != 1 || !strings.Contains(w.Body.String(), "第二题") {
		t.Fatalf("expected second practice question, got %d", w.Code)
	}
	if strings.Contains(w.Body.String(), `name="direction" value="next"`) {
		t.Fatal("last question rendered next navigation")
	}
	w = shared.Serve(f, f.RequestAs(t, 2, http.MethodGet, "/practice/101", nil))
	if w.Code != http.StatusNotFound {
		t.Fatalf("owner validation failed: got %d", w.Code)
	}
	p.Completed = true
	w = shared.Serve(f, f.Request(http.MethodGet, "/practice/101", nil))
	if w.Code != http.StatusGone {
		t.Fatalf("completed practice: got %d", w.Code)
	}

	now := time.Now()
	for id := 1; id <= 6; id++ {
		record := &model.PracticeRecord{UserID: 1, PracticeID: 200 + id, TotalQuestions: 1, CorrectCount: 1, StartTime: now, SubmitTime: now.Add(time.Duration(id) * time.Second)}
		if err := f.Server.PM.Persist(record); err != nil {
			t.Fatal(err)
		}
	}
	w = shared.Serve(f, f.Request(http.MethodGet, "/practice/history", nil))
	if w.Code != http.StatusOK || !strings.Contains(w.Body.String(), "206") || !strings.Contains(w.Body.String(), "page=2") {
		t.Fatalf("unexpected history page: %d", w.Code)
	}
	w = shared.Serve(f, f.Request(http.MethodGet, "/practice/history?page=0", nil))
	if w.Code != http.StatusBadRequest {
		t.Fatalf("invalid history page: got %d", w.Code)
	}
}

// TestPracticeGetPersistedResult 验证持久化套题结果页面 Handler。
func TestPracticeGetPersistedResult(t *testing.T) {
	f := shared.NewFixture(t)
	now := time.Now()
	record := &model.PracticeRecord{UserID: 1, PracticeID: 301, TotalQuestions: 1, CorrectCount: 1, StartTime: now, SubmitTime: now, Answers: []model.AnswerRecord{{QuestionID: 1, Answer: 1, Answered: true, Correct: true}}}
	if err := f.Server.PM.Persist(record); err != nil {
		t.Fatal(err)
	}
	w := shared.Serve(f, f.Request(http.MethodGet, "/practice/301/result", nil))
	if w.Code != http.StatusOK || !strings.Contains(w.Body.String(), "第一题解析") {
		t.Fatalf("unexpected result page: %d", w.Code)
	}
}
