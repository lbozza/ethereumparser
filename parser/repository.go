package parser

import (
	"eth/domain"
	"eth/repository"
)

type Repository interface {
	Store(address string, transaction domain.Transaction)
	GetTransactions(address string) []domain.Transaction
}

type repo struct {
	storage *repository.Memory
}

func NewRepository(storage *repository.Memory) *repo {
	return &repo{
		storage: storage,
	}
}
