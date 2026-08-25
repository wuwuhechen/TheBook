package manager

import "TheBook/model/structs"

type QuestionProgressManager interface {
	FindByUserID(userID uint) (*structs.QuestionProgress, error)

	Upsert(progress *structs.QuestionProgress) error

	Persist() error
}
