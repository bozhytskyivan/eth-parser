package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEthClient_GetCurrentBlock(t *testing.T) {
	cl := NewEthClient()

	currentBlock, err := cl.GetCurrentBlock(DefaultID)
	assert.NoError(t, err)
	assert.True(t, currentBlock > 0)
}

func TestEthClient_GetBlockByNumber(t *testing.T) {
	cl := NewEthClient()

	const (
		blockHash               = "0xb3b20624f8f0f86eb50dd04688409e5cea4bd02d700bf6e79e9384d47d6a5a35"
		blockNumber             = 0x5bad55
		blockTransactionsAmount = 80
	)

	block, err := cl.GetBlockByNumber(DefaultID, blockNumber)
	assert.NoError(t, err)
	assert.Equal(t, blockHash, block.Hash)
	assert.Equal(t, blockTransactionsAmount, len(block.Transactions))
}
