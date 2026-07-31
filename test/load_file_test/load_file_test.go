package load_file_test

import (
	"TheBook/service"
	"fmt"
	"testing"
)

func TestLoadQuestions(t *testing.T) {
	qs, err := service.LoadQuestions("../../database/data.json")
	if err != nil {
		t.Fatalf("Failed to load questions: %v", err)
	}

	db := qs.DB.GetAllQuestionsSorted()

	for _, question := range db {
		fmt.Printf("ID: %d, Question: %s, Choices: %v\n", question.ID, question.Question, question.Choices)
	}
}

func BenchmarkLoadQuestions(b *testing.B) {

	for i := 0; i < b.N; i++ {

		_, err := service.LoadQuestions("../../database/data.json")

		if err != nil {
			b.Fatal(err)
		}
	}
}
