package structs

import (
	"strings"
	"testing"
	"time"
)

func TestQuestionPageDataMapsQuestionAndSetters(t *testing.T) {
	question := &Question{ID: 3, Category: "Go", Question: "question", Choices: []string{"A", "B", "C", "D"}, Explanation: "because"}
	data := NewQuestionData(question, 5)
	if data.RealID != 3 || !data.HasLastID || !data.HasNextID || data.LastID != 2 || data.NextID != 4 || data.OptionD != "D" {
		t.Fatalf("unexpected page data: %+v", data)
	}
	data.SetLastID(30)
	data.SetNextID(50, 8)
	data.SetHasLastID(false)
	data.SetHasNextID(false)
	data.SetDuration(95 * time.Second)
	data.SetPracticeID(7)
	data.SetRandomSessionID(8)
	data.SetUserID(9)
	data.SetID(2)
	if data.LastID != 30 || data.NextID != 50 || data.TotalQuestions != 8 || data.HasLastID || data.HasNextID || data.Duration != 95 || data.PracticeID != 7 || data.RandomSessionID != 8 || data.UserID != 9 || data.ID != 2 {
		t.Fatalf("setters did not update page data: %+v", data)
	}
	if !strings.Contains(data.String(), "编号：2") || !strings.Contains(data.String(), "问题：question") {
		t.Fatalf("unexpected string representation: %s", data.String())
	}
}
