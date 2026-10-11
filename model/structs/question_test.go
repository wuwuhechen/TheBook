package structs

import (
	"strings"
	"testing"
)

// TestQuestionCheckAnswerAndString 验证题目判题和字符串表示。
func TestQuestionCheckAnswerAndString(t *testing.T) {
	question := Question{
		ID: 7, Category: "Go", Question: "Go 的包名？",
		Choices: []string{"A", "B", "C", "D"}, Answer: 2, Explanation: "第三项正确",
	}

	if !question.CheckAnswer(2) {
		t.Fatal("expected correct choice to pass")
	}
	if question.CheckAnswer(1) {
		t.Fatal("expected incorrect choice to fail")
	}
	for _, expected := range []string{"编号: 7", "分类: Go", "Go 的包名？", "答案: 2", "第三项正确"} {
		if !strings.Contains(question.String(), expected) {
			t.Fatalf("question string does not contain %q: %s", expected, question.String())
		}
	}
}
