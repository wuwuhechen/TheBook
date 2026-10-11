package wrong_question_test

import (
	"TheBook/test/gin_handler_test/shared"
	"net/http"
	"strings"
	"testing"
	"time"
)

// TestWrongQuestionGetPage 验证错题本页面 Handler。
func TestWrongQuestionGetPage(t *testing.T) {
	f := shared.NewFixture(t)
	if err := f.Server.WQ.RecordWrongQuestion(1, 1, time.Now()); err != nil {
		t.Fatal(err)
	}
	w := shared.Serve(f, f.Request(http.MethodGet, "/practice/wrong-questions", nil))
	if w.Code != http.StatusOK || !strings.Contains(w.Body.String(), "第一题") {
		t.Fatalf("wrong-question page failed: %d", w.Code)
	}
	w = shared.Serve(f, f.Request(http.MethodGet, "/practice/wrong-questions?page=bad", nil))
	if w.Code != http.StatusBadRequest {
		t.Fatalf("invalid wrong-question page: got %d", w.Code)
	}
}
