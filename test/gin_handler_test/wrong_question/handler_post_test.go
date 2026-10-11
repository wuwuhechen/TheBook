package wrong_question_test

import (
	"TheBook/model"
	"TheBook/test/gin_handler_test/shared"
	"net/http"
	"strings"
	"testing"
	"time"
)

// TestWrongQuestionPostPractice 验证错题重练创建 Handler。
func TestWrongQuestionPostPractice(t *testing.T) {
	f := shared.NewFixture(t)
	w := shared.Serve(f, f.Request(http.MethodPost, "/wrong-questions/practice", shared.JSONBody(t, model.Request{})))
	if w.Code != http.StatusBadRequest {
		t.Fatalf("empty selection: got %d", w.Code)
	}
	w = shared.Serve(f, f.Request(http.MethodPost, "/wrong-questions/practice", shared.JSONBody(t, model.Request{QuestionIDs: []int{1, 1}})))
	if w.Code != http.StatusBadRequest {
		t.Fatalf("duplicate selection: got %d", w.Code)
	}
	w = shared.Serve(f, f.Request(http.MethodPost, "/wrong-questions/practice", shared.JSONBody(t, model.Request{QuestionIDs: []int{1}})))
	if w.Code != http.StatusForbidden {
		t.Fatalf("unowned selection: got %d", w.Code)
	}
	if err := f.Server.WQ.RecordWrongQuestion(1, 1, time.Now()); err != nil {
		t.Fatal(err)
	}
	w = shared.Serve(f, f.Request(http.MethodPost, "/wrong-questions/practice", shared.JSONBody(t, model.Request{QuestionIDs: []int{1}})))
	if w.Code != http.StatusFound || !strings.HasPrefix(w.Header().Get("Location"), "/practice/") {
		t.Fatalf("wrong practice creation failed: %d", w.Code)
	}
}
