package structs

import (
	"fmt"
	"testing"
	"time"
)

type practiceTestQuestionManager struct {
	questions map[int]Question
	ids       []int
}

func (m *practiceTestQuestionManager) GetQuestion(id int) (*Question, error) {
	question, ok := m.questions[id]
	if !ok {
		return nil, fmt.Errorf("question %d not found", id)
	}
	return &question, nil
}
func (m *practiceTestQuestionManager) GetAllQuestionsSorted() []Question { return m.GetAllQuestions() }
func (m *practiceTestQuestionManager) GetAllQuestions() []Question {
	questions := make([]Question, 0, len(m.ids))
	for _, id := range m.ids {
		questions = append(questions, m.questions[id])
	}
	return questions
}
func (m *practiceTestQuestionManager) GetALLQuestionIDs() []int { return append([]int(nil), m.ids...) }
func (m *practiceTestQuestionManager) GetRandomQuestionID() (*Question, error) {
	return m.GetQuestion(m.ids[0])
}
func (m *practiceTestQuestionManager) GetTotalCount() int { return len(m.ids) }

func newPracticeTestManager() *practiceTestQuestionManager {
	return &practiceTestQuestionManager{
		ids: []int{10, 20, 30},
		questions: map[int]Question{
			10: {ID: 10, Answer: 0, Explanation: "ten"},
			20: {ID: 20, Answer: 1, Explanation: "twenty"},
			30: {ID: 30, Answer: 2, Explanation: "thirty"},
		},
	}
}

func TestNewPracticeHasSafeDefaults(t *testing.T) {
	practice := (&Practice{}).NewPractice()
	if practice.Answers == nil || practice.Questions == nil {
		t.Fatal("new practice must initialize collections")
	}
	if practice.Duration != time.Hour || practice.Completed || practice.CurrentIndex != 0 {
		t.Fatalf("unexpected defaults: %+v", practice)
	}
}

func TestGenerateExamUsesUniqueQuestionsAndRequestedSize(t *testing.T) {
	manager := newPracticeTestManager()
	practice := (&Practice{}).GenerateExam(manager, 2)
	if practice.PracticeSize != 2 || practice.TotalQuestions != 2 || len(practice.Questions) != 2 {
		t.Fatalf("unexpected generated practice: %+v", practice)
	}
	seen := make(map[int]bool)
	for _, id := range practice.Questions {
		if seen[id] || manager.questions[id].ID == 0 {
			t.Fatalf("invalid or duplicate question ID %d", id)
		}
		seen[id] = true
	}
	if practice.ID < 100000000 || practice.ID > 999999999 {
		t.Fatalf("expected a nine digit practice ID, got %d", practice.ID)
	}

	all := (&Practice{}).GenerateExam(manager, 99)
	if len(all.Questions) != 3 || all.TotalQuestions != 3 {
		t.Fatalf("expected size to be capped by question count: %+v", all)
	}
}

func TestPracticeNavigationAndReset(t *testing.T) {
	practice := &Practice{Questions: []int{10, 20, 30}, TotalQuestions: 3, CurrentIndex: 0, Answers: map[int]int{10: 1}, Completed: true}
	if practice.GetCurrentQuestionID() != 10 || practice.LastQuestionID() != 0 || practice.NextQuestionID() != 20 {
		t.Fatalf("unexpected first-question navigation")
	}
	practice.CurrentIndex = 1
	if practice.LastQuestionID() != 10 || practice.GetCurrentQuestionID() != 20 || practice.NextQuestionID() != 30 {
		t.Fatalf("unexpected middle-question navigation")
	}
	practice.CurrentIndex = 2
	if practice.NextQuestionID() != practice.TotalQuestions {
		t.Fatal("last question must not expose a next question ID")
	}
	practice.Reset()
	if practice.CurrentIndex != 0 || len(practice.Answers) != 0 || practice.Completed || practice.StartTime.IsZero() {
		t.Fatalf("reset did not restore a new practice state: %+v", practice)
	}
}

func TestPracticeCheckPracticeCountsAnswers(t *testing.T) {
	manager := newPracticeTestManager()
	practice := &Practice{ID: 99, Questions: []int{10, 20, 30}, Answers: map[int]int{10: 0, 20: 0}}
	result := practice.CheckPractice(manager)
	if result.PracticeID != 99 || result.Total != 3 || result.CorrectCount != 1 || result.WrongCount != 2 {
		t.Fatalf("unexpected practice result: %+v", result)
	}
	if len(result.Details) != 3 || !result.Details[0].Correct || result.Details[1].Correct || result.Details[2].Correct {
		t.Fatalf("unexpected result details: %+v", result.Details)
	}
}

func TestGenerateExamByQuestionIDsCopiesInput(t *testing.T) {
	ids := []int{10, 20, 30}
	practice := (&Practice{}).GenerateExamByQuestionIDs(ids)
	ids[0] = 999
	for _, id := range practice.Questions {
		if id == 999 {
			t.Fatal("practice must not share caller's question ID slice")
		}
	}
	if practice.GetDuration() != time.Hour {
		t.Fatal("expected default duration")
	}
}
