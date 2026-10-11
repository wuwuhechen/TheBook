package question_test

import (
	"TheBook/test/gin_handler_test/shared"
	"net/http"
	"strings"
	"testing"
)

// TestQuestionGetAndCurrent 验证单题页面和继续答题 Handler。
func TestQuestionGetAndCurrent(t *testing.T) {
	f := shared.NewFixture(t)
	for name, target := range map[string]string{"invalid": "/question?question_id=bad", "missing": "/question?question_id=99", "valid": "/question?question_id=1"} {
		t.Run(name, func(t *testing.T) {
			w := shared.Serve(f, f.Request(http.MethodGet, target, nil))
			if name == "valid" {
				if w.Code != http.StatusOK || !strings.Contains(w.Body.String(), "第一题") {
					t.Fatalf("expected question page, got %d", w.Code)
				}
			} else if w.Code != http.StatusBadRequest && w.Code != http.StatusNotFound {
				t.Fatalf("expected client error, got %d", w.Code)
			}
		})
	}
	progress, err := f.Server.QS.FindByUserID(1)
	if err != nil {
		t.Fatal(err)
	}
	progress.CurrentQuestionID = 999
	if err := f.Server.QS.Upsert(progress); err != nil {
		t.Fatal(err)
	}
	w := shared.Serve(f, f.Request(http.MethodPost, "/question/current", nil))
	if w.Code != http.StatusFound || w.Header().Get("Location") != "/question?question_id=1" {
		t.Fatalf("unexpected current redirect: %d %q", w.Code, w.Header().Get("Location"))
	}
}
