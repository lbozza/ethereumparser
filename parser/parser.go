package parser

import (
	"eth/domain"
	"sync"
)

type Parser interface {
	// GetCurrentBlock last parsed block
	GetCurrentBlock() int

	// Subscribe add address to observer
	Subscribe(address string) bool

	// GetTransactions list of inbound or outbound transactions for an address
	GetTransactions(address string) []domain.Transaction
}

type EthereumParser struct {
	processor  *Processor
	repo       Repository
	wg         sync.WaitGroup
	intialized bool
}

func NewParser(url string) *EthereumParser {
	memoryStorage := NewMemory()
	repo := NewRepository(memoryStorage)

	parser := &EthereumParser{
		processor:  NewProcessor(url, repo.storage),
		repo:       repo.storage,
		wg:         sync.WaitGroup{},
		intialized: false,
	}

	return parser
}

func (parser *EthereumParser) GetCurrentBlock() int {
	return int(parser.processor.latestProcessedBlock)
}

func (parser *EthereumParser) Subscribe(address string) bool {
	parser.repo.Subscribe(address)
	return true
}

func (parser *EthereumParser) GetTransactions(address string) []domain.Transaction {
	return parser.repo.GetTransactions(address)
}
