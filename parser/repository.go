package parser

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

func (m *Memory) Subscribe(address string) bool {
	m.lock.Lock()
	defer m.lock.Unlock()
	if _, ok := m.data[address]; ok {
		return false
	}
	m.data[address] = []domain.Transaction{}
	return true
}

func (m *Memory) IsSubscribed(address string) bool {
	m.lock.RLock()
	defer m.lock.RUnlock()
	_, ok := m.data[address]
	return ok
}

type Repository interface {
	Store(address string, transaction domain.Transaction)
	GetTransactions(address string) []domain.Transaction
	Subscribe(address string) bool
	IsSubscribed(address string) bool
}

type repo struct {
	storage *Memory
}

func NewRepository(storage *Memory) *repo {
	return &repo{
		storage: storage,
	}
}
