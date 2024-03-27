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
	blockSeen            map[int64]time.Time
	ticker               time.Ticker
	tickerCleanBlocks    time.Ticker
	repository           Repository
}

func NewProcessor(url string, repository Repository) *Processor {
	return &Processor{
		lock:                 sync.RWMutex{},
		ETHClient:            *NewETHClient(url),
		latestProcessedBlock: 0,
		subscribers:          make(map[string]bool),
		blocks:               make(chan int64, 1000),
		blockSeen:            make(map[int64]time.Time),
		ticker:               *time.NewTicker(time.Second * 5),
		tickerCleanBlocks:    *time.NewTicker(time.Second * 60),
		repository:           repository,
	}
}

func (p *Processor) Start() {
	fmt.Println("Starting")
	for range p.ticker.C {
		block, err := p.ETHClient.GetBlockNumber()
		if err != nil {
			p.prettyLog("Error getting block number: %v", err)
			continue // Continue to the next iteration of the loop
		}

		var found bool
		blockNum := FormatHexToInt(block.Result)
		for k := range p.blockSeen {
			if k == blockNum {
				found = true
			}
		}
		if !found {
			p.prettyLog("Adding block to channel: %d", blockNum)
			p.blocks <- FormatHexToInt(block.Result)
		}

	}
}

func (p *Processor) ProcessBlock() {
	for {
		select {
		case b := <-p.blocks:
			blockResponse, err := p.ETHClient.GetBlockByNumber(FormatIntToHex(b))
			if err != nil {
				fmt.Printf("Error trying to get block %d information: %d\n", b, err)
				continue
			}

			p.prettyLog("Processing block: %d", b)
			p.blockSeen[b] = time.Now()

			p.latestProcessedBlock = b

			for _, tx := range blockResponse.Result.Transactions {
				if p.subscribers[tx.To] {
					p.prettyLog("Found Transaction in block %d to address: %s ", FormatHexToInt(tx.BlockNumber), tx.To)
					p.repository.Store(tx.To, tx)
				} else if p.subscribers[tx.From] {
					p.prettyLog("Found Transaction in block %d from address: %s ", FormatHexToInt(tx.BlockNumber), tx.From)
					p.repository.Store(tx.From, tx)
				}
			}
		}
	}
}

func (p *Processor) cleanProcessedBlocks() {
	p.lock.Lock()
	defer p.lock.Unlock()
	for range p.tickerCleanBlocks.C {
		for k, v := range p.blockSeen {
			if v.Before(time.Now().Add(-5 * time.Minute)) {
				p.prettyLog("Deleting Processed Block %d From seen queue", k)
				delete(p.blockSeen, k)
			}
		}
	}
}

func (p *Processor) prettyLog(msg string, args ...interface{}) {
	fmt.Printf("[Processor] "+msg+"\n", args...)
}
