package repository

import (
	"eth/domain"
	"sync"
)

type Memory struct {
	lock sync.RWMutex
	data map[string][]domain.Transaction
}

func NewMemory() *Memory {
	return &Memory{
		sync.RWMutex{},
		make(map[string][]domain.Transaction),
	}
}

func (m *Memory) Store(address string, tx domain.Transaction) {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.data[address] = append(m.data[address], tx)
}

func (m *Memory) GetTransactions(address string) []domain.Transaction {
	m.lock.RLock()
	defer m.lock.RUnlock()
	return m.data[address]
}
