package manager

import "TheBook/model/structs"

type RandomSessionManager interface {
	Create(session *structs.RandomSession) error

	FindByID(sessionID int) (*structs.RandomSession, error)

	Save(session *structs.RandomSession) error

	Delete(sessionID int) error
}
