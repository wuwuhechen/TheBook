package model

type RandomSessionManager interface {
	Create(session *RandomSession) error

	FindByID(sessionID int) (*RandomSession, error)

	Save(session *RandomSession) error

	Delete(sessionID int) error
}
