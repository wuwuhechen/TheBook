package bank

import (
	"TheBook/model/manager"
	"errors"
	"sync"
)

type RandomSessionBank struct {
	randomSessions map[int]*RandomSession

	mu sync.RWMutex
}

var _ manager.RandomSessionManager = (*RandomSessionBank)(nil)

func NewRandomSessionBank() *RandomSessionBank {
	return &RandomSessionBank{
		randomSessions: make(map[int]*RandomSession),
	}
}

func (b *RandomSessionBank) Create(session *RandomSession) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.randomSessions[session.ID] = session
	return nil
}

func (b *RandomSessionBank) FindByID(sessionID int) (*RandomSession, error) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	session, ok := b.randomSessions[sessionID]
	if !ok {
		return nil, errors.New("random session not found")
	}
	return session, nil
}

func (b *RandomSessionBank) Save(session *RandomSession) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.randomSessions[session.ID] = session
	return nil
}

func (b *RandomSessionBank) Delete(sessionID int) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	delete(b.randomSessions, sessionID)
	return nil
}
