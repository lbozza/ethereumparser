package parser

import (
	"eth/domain"
	"eth/repository"
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
	intialized bool
}

func NewParser(url string) *EthereumParser {
	memoryStorage := repository.NewMemory()
	repo := NewRepository(memoryStorage)

	parser := &EthereumParser{
		processor:  NewProcessor(url, repo.storage),
		repo:       repo.storage,
		intialized: false,
	}

	return parser
}

func (parser *EthereumParser) run() {

	go parser.processor.Start()
	go parser.processor.cleanProcessedBlocks()
	parser.processor.ProcessBlock()

}

func (parser *EthereumParser) GetCurrentBlock() int {
	return int(parser.processor.latestProcessedBlock)
}

func (parser *EthereumParser) Subscribe(address string) bool {
	if !parser.intialized {
		parser.intialized = true
		go parser.run()
	}
	parser.processor.subscribers[address] = true
	return true
}

func (parser *EthereumParser) GetTransactions(address string) []domain.Transaction {
	return parser.repo.GetTransactions(address)
}
