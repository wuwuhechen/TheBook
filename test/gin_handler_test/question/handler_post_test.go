package question_test

import (
	"TheBook/model"
	"TheBook/test/gin_handler_test/shared"
	"net/http"
	"strings"
	"testing"
)

// TestQuestionPosts 验证单题跳转和判题 Handler。
func TestQuestionPosts(t *testing.T) {
	f := shared.NewFixture(t)
	w := shared.Serve(f, f.Request(http.MethodPost, "/question/request", []byte(`{`)))
	if w.Code != http.StatusBadRequest {
		t.Fatalf("invalid request: got %d", w.Code)
	}
	w = shared.Serve(f, f.Request(http.MethodPost, "/question/request", shared.JSONBody(t, model.Request{QuestionID: 99})))
	if w.Code != http.StatusNotFound {
		t.Fatalf("missing request: got %d", w.Code)
	}
	w = shared.Serve(f, f.Request(http.MethodPost, "/question/request", shared.JSONBody(t, model.Request{QuestionID: 2})))
	if w.Code != http.StatusFound || w.Header().Get("Location") != "/question?question_id=2" {
		t.Fatalf("unexpected redirect: %d %q", w.Code, w.Header().Get("Location"))
	}
	w = shared.Serve(f, f.Request(http.MethodPost, "/question/check_answer", []byte(`{`)))
	if w.Code != http.StatusBadRequest {
		t.Fatalf("invalid answer: got %d", w.Code)
	}
	w = shared.Serve(f, f.Request(http.MethodPost, "/question/check_answer", shared.JSONBody(t, model.Request{UserID: 1, QuestionID: 1, Choice: 0})))
	if w.Code != http.StatusOK || !strings.Contains(w.Body.String(), `"correct":false`) {
		t.Fatalf("expected incorrect response, got %d: %s", w.Code, w.Body.String())
	}
	wrong, err := f.Server.WQ.ListWrongQuestionsByUser(1, 10, 0)
	if err != nil || len(wrong) != 1 {
		t.Fatalf("wrong question not persisted: %+v, %v", wrong, err)
	}
}
