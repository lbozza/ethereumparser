package parser

import (
	"fmt"
	"sync"
	"time"
)

type Processor struct {
	lock                 sync.RWMutex
	ETHClient            ETHClient
	latestProcessedBlock int64
	subscribers          map[string]bool
	blocks               chan int64
	ticker               time.Ticker
	repository           Repository
}

func NewProcessor(url string, repository Repository) *Processor {
	p := &Processor{
		lock:                 sync.RWMutex{},
		ETHClient:            *NewETHClient(url),
		latestProcessedBlock: 0,
		subscribers:          make(map[string]bool),
		blocks:               make(chan int64, 1000),
		ticker:               *time.NewTicker(time.Second * 5),
		repository:           repository,
	}

	go p.Start()
	go p.ProcessBlock()

	return p
}

func (p *Processor) Start() {
	for range p.ticker.C {
		block, err := p.ETHClient.GetBlockNumber()
		if err != nil {
			p.prettyLog("Error getting block number: %v", err)
			continue
		}

		blockNum := FormatHexToInt(block.Result)

		if blockNum > p.getCurrentBlock() {
			p.blocks <- blockNum

		}

	}
}

func (p *Processor) ProcessBlock() {
	for {
		select {
		case b := <-p.blocks:

			if p.latestProcessedBlock == 0 {
				p.updateLastProcessedBlock(b - 1)
			}

			if b > p.getCurrentBlock() {
				for i := p.getCurrentBlock() + 1; i <= b; i++ {
					blockResponse, err := p.ETHClient.GetBlockByNumber(FormatIntToHex(b))
					if err != nil {
						fmt.Printf("Error trying to get block %d information: %d\n", b, err)
						continue
					}

					p.prettyLog("Processing block: %d", b)

					p.updateLastProcessedBlock(FormatHexToInt(blockResponse.Result.Number))

					for _, tx := range blockResponse.Result.Transactions {
						if p.repository.IsSubscribed(tx.To) {
							p.prettyLog("Found Transaction in block %d to address: %s ", FormatHexToInt(tx.BlockNumber), tx.To)
							p.repository.Store(tx.To, tx)
						} else if p.repository.IsSubscribed(tx.From) {
							p.prettyLog("Found Transaction in block %d from address: %s ", FormatHexToInt(tx.BlockNumber), tx.From)
							p.repository.Store(tx.From, tx)
						}
					}
				}

			}

		}
	}
}

func (p *Processor) updateLastProcessedBlock(block int64) {
	p.lock.Lock()
	defer p.lock.Unlock()

	p.latestProcessedBlock = block
}

func (p *Processor) getCurrentBlock() int64 {
	p.lock.RLock()
	defer p.lock.RUnlock()

	return p.latestProcessedBlock
}

func (p *Processor) prettyLog(msg string, args ...interface{}) {
	fmt.Printf("[Processor] "+msg+"\n", args...)
}
