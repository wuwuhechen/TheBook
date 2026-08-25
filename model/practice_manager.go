package model

type PracticeManager interface {
	Create(practice *Practice) error

	FindByID(id int) (*Practice, error)

	Save(practice *Practice) error

	Delete(practice *Practice) error
}

var _ PracticeManager = (*PracticeBank)(nil)
