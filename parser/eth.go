package parser

import (
	"encoding/json"
	"eth/domain"
	"net/http"
	"strings"
	"sync"
)

const (
	rpcVersion       = "2.0"
	reqEncoding      = "application/json"
	ethBlockNumber   = "eth_blockNumber"
	ethBlockByNumber = "eth_getBlockByNumber"
)

type ETHClient struct {
	url      string
	sequence int
	mu       sync.Mutex
}

func NewETHClient(url string) *ETHClient {
	return &ETHClient{
		url:      url,
		sequence: 0,
	}
}

func (c *ETHClient) GetBlockNumber() (*domain.BlockNumber, error) {
	response, err := c.doRequest(ethBlockNumber, []interface{}{})
	if err != nil {
		return nil, err
	}

	block := &domain.BlockNumber{}
	err = json.NewDecoder(response.Body).Decode(block)
	if err != nil {
		return nil, err
	}
	defer closeResponse(response)
	return block, err
}

func (c *ETHClient) GetBlockByNumber(blockNumber string) (*domain.Block, error) {
	response, err := c.doRequest(ethBlockByNumber, []interface{}{blockNumber, true})
	if err != nil {
		return nil, err
	}

	block := &domain.Block{}
	err = json.NewDecoder(response.Body).Decode(block)
	if err != nil {
		return nil, err
	}

	defer closeResponse(response)
	return block, err
}

func (c *ETHClient) doRequest(method string, params []interface{}) (*http.Response, error) {
	c.updateSequence()
	r := domain.Request{
		JsonRPC: rpcVersion,
		Method:  method,
		Params:  params,
		ID:      c.sequence,
	}

	m, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}

	response, err := http.Post(c.url, reqEncoding, strings.NewReader(string(m)))
	if err != nil {
		return nil, err
	}

	return response, err
}

func (c *ETHClient) updateSequence() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.sequence++
}

func closeResponse(response *http.Response) {
	if response != nil {
		response.Body.Close()
	}
}
