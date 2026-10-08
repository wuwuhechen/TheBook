package structs

import "testing"

func TestRandomSessionCurrentQuestionIDBounds(t *testing.T) {
	session := &RandomSession{Questions: []int{10, 20}}
	if session.CurrentQuestionID() != 10 {
		t.Fatal("expected first question")
	}
	session.CurrentIndex = 1
	if session.CurrentQuestionID() != 20 {
		t.Fatal("expected second question")
	}
	for _, index := range []int{-1, 2} {
		session.CurrentIndex = index
		if session.CurrentQuestionID() != 0 {
			t.Fatalf("expected 0 for index %d", index)
		}
	}
}

func TestNewRandomSessionIncludesEveryQuestion(t *testing.T) {
	manager := newPracticeTestManager()
	session := NewRandomSession(manager)
	if len(session.Questions) != 3 || session.CurrentIndex != 0 {
		t.Fatalf("unexpected session: %+v", session)
	}
	seen := make(map[int]bool)
	for _, id := range session.Questions {
		seen[id] = true
	}
	for _, id := range []int{10, 20, 30} {
		if !seen[id] {
			t.Fatalf("session is missing question %d", id)
		}
	}
}
