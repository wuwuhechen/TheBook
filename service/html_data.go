package service

type QuestionPageData struct {
	ID             int
	Category       string
	Question       string
	OptionA        string
	OptionB        string
	OptionC        string
	OptionD        string
	TotalQuestions int
}

func NewQuestionData(question *Question, totalQuestions int) *QuestionPageData {
	return &QuestionPageData{
		ID:             question.ID,
		Category:       question.Category,
		Question:       question.Question,
		OptionA:        question.Choices[0],
		OptionB:        question.Choices[1],
		OptionC:        question.Choices[2],
		OptionD:        question.Choices[3],
		TotalQuestions: totalQuestions,
	}
}
