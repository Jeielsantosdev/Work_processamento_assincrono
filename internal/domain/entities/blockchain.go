package entities

import (
	"time"
)

// Block representa um bloco na blockchain privada
type Block struct {
	Index        int64
	Timestamp    time.Time
	Transactions []Transaction
	PreviousHash string
	Hash         string
	Nonce        int64
	Difficulty   int
}

// NewBlock cria um novo bloco
func NewBlock(index int64, transactions []Transaction, previousHash string) *Block {
	return &Block{
		Index:        index,
		Timestamp:    time.Now(),
		Transactions: transactions,
		PreviousHash: previousHash,
		Hash:         "",
		Nonce:        0,
		Difficulty:   2,
	}
}

// BlockRecord representa um registro auditável de uma transação no blockchain
type BlockRecord struct {
	ID              string
	TransactionID   string
	BlockIndex      int64
	BlockHash       string
	PreviousHash    string
	Timestamp       time.Time
	Status          string
	TransactionData map[string]interface{}
}

// NewBlockRecord cria um novo registro de blockchain
func NewBlockRecord(txID string, blockIndex int64, blockHash, previousHash string, txData map[string]interface{}) *BlockRecord {
	return &BlockRecord{
		ID:              txID + "-" + string(rune(blockIndex)),
		TransactionID:   txID,
		BlockIndex:      blockIndex,
		BlockHash:       blockHash,
		PreviousHash:    previousHash,
		Timestamp:       time.Now(),
		Status:          "RECORDED",
		TransactionData: txData,
	}
}
