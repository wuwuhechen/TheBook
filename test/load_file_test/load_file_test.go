package load_file_test

import (
	"TheBook/service"
	"TheBook/utils"
	"fmt"
	"testing"
)

func TestDataInit(t *testing.T) {
	rootPath, err := utils.FindProjectRoot()
	if err != nil {
		t.Fatalf("Failed to find project root: %v", err)
	}

	qs, err := service.DataInit(fmt.Sprintf("%s/database/data.json", rootPath))
	if err != nil {
		t.Fatalf("Failed to initialize system: %v", err)
	}

	db := qs.DB.GetAllQuestionsSorted()

	for _, question := range db {
		fmt.Printf("ID: %d, Question: %s, Explanation: %s\n", question.ID, question.Question, question.Explanation)
	}
}

func BenchmarkLoadQuestions(b *testing.B) {

	rootPath, err := utils.FindProjectRoot()
	if err != nil {
		b.Fatalf("Failed to find project root: %v", err)
	}

	for i := 0; i < b.N; i++ {

		_, err := service.LoadQuestions(fmt.Sprintf("%s/database/data.json", rootPath))

		if err != nil {
			b.Fatal(err)
		}
	}
}
