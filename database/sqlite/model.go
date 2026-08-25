package main

import "time"

type Question struct {
	ID uint `gorm:"primaryKey" json:"id"`

	Category string `gorm:"type:text;not null" json:"category"`
	Question string `gorm:"type:text;not null" json:"question"`

	OptionA string `gorm:"type:text;not null" json:"option_a"`
	OptionB string `gorm:"type:text;not null" json:"option_b"`
	OptionC string `gorm:"type:text;not null" json:"option_c"`
	OptionD string `gorm:"type:text;not null" json:"option_d"`

	Answer      int    `gorm:"type:integer;not null" json:"answer" `
	Explanation string `gorm:"type:text;not null" json:"explanation"`
}

type User struct {
	UserID   uint   `gorm:"primaryKey" json:"user_id"`
	Username string `gorm:"type:text;not null" json:"username"`
	Password string `gorm:"type:text;not null" json:"password"`
	Nickname string `gorm:"type:text;not null" json:"nickname"`
	Role     string `gorm:"type:text;not null" json:"role"`

	QuestionProgress QuestionProgress `gorm:"foreignKey:UserID"`
	PracticeRecords  []PracticeRecord `gorm:"foreignKey:UserID"`
}

type PracticeAnswer struct {
	ID               uint `gorm:"primaryKey" json:"id"`
	PracticeRecordID uint `gorm:"not null;uniqueIndex:idx_practice_answer" json:"practice_record_id"`
	QuestionID       uint `gorm:"not null;uniqueIndex:idx_practice_answer" json:"question_id"`
	Answer           int  `gorm:"not null" json:"answer"`
	Answered         bool `gorm:"not null" json:"answered"`
	Correct          bool `gorm:"not null" json:"correct"`
}

type PracticeRecord struct {
	ID     uint `gorm:"primaryKey" json:"id"`
	UserID uint `gorm:"not null" json:"user_id"`

	PracticeID     uint `gorm:"not null" json:"practice_id"`
	TotalQuestions int  `gorm:"not null" json:"total_questions"`

	CorrectCount int `gorm:"not null" json:"correct_count"`
	WrongCount   int `gorm:"not null" json:"wrong_count"`

	StartTime  time.Time `gorm:"not null" json:"start_time"`
	SubmitTime time.Time `gorm:"not null" json:"submit_time"`

	Answers []PracticeAnswer `gorm:"foreignKey:PracticeRecordID"`
}

type QuestionProgress struct {
	ID                uint      `gorm:"primaryKey" json:"id"`
	UserID            uint      `gorm:"not null;uniqueIndex" json:"user_id"`
	CurrentQuestionID uint      `gorm:"not null" json:"current_question_id"`
	UpdatedAt         time.Time `gorm:"not null" json:"updated_at"`
}
