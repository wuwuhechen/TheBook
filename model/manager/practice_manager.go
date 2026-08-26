package manager

import "TheBook/model/structs"

type PracticeManager interface {
	Create(practice *structs.Practice) error

	FindByID(id int) (*structs.Practice, error)

	Save(practice *structs.Practice) error

	Delete(practice *structs.Practice) error

	Persist(record *structs.PracticeRecord) error
}
