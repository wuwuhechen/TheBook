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

// NewRandomSessionBank 创建空的随机答题会话仓库。
func NewRandomSessionBank() *RandomSessionBank {
	return &RandomSessionBank{
		randomSessions: make(map[int]*RandomSession),
	}
}

// Create 保存新的随机答题会话。
func (b *RandomSessionBank) Create(session *RandomSession) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.randomSessions[session.ID] = session
	return nil
}

// FindByID 按会话 ID 查询随机答题会话。
func (b *RandomSessionBank) FindByID(sessionID int) (*RandomSession, error) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	session, ok := b.randomSessions[sessionID]
	if !ok {
		return nil, errors.New("random session not found")
	}
	return session, nil
}

// Save 更新随机答题会话的当前状态。
func (b *RandomSessionBank) Save(session *RandomSession) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.randomSessions[session.ID] = session
	return nil
}

// Delete 删除指定随机答题会话。
func (b *RandomSessionBank) Delete(sessionID int) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	delete(b.randomSessions, sessionID)
	return nil
}
