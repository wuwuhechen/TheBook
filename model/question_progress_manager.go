package model

type QuestionProgressManager interface {
	FindByUserID(userID uint) (*QuestionProgress, error)

	Upsert(progress *QuestionProgress) error

	Persist() error
}
