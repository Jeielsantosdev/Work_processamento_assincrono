package blockchain

import (
	"context"
	"crypto/sha256"
	"fmt"
	"time"

	"github.com/Jeielsantosdev/Work_processamento_assincrono/internal/domain/entities"
	"github.com/Jeielsantosdev/Work_processamento_assincrono/internal/domain/interfaces"
)

// SimpleBlockchainService implementa um serviço de blockchain simples
// Em produção, seria integrado com um blockchain real
type SimpleBlockchainService struct {
	blockchainRepo interfaces.BlockchainRepository
	lastBlockHash  string
	blockIndex     int64
}

// NewSimpleBlockchainService cria uma nova instância do serviço de blockchain
func NewSimpleBlockchainService(blockchainRepo interfaces.BlockchainRepository) *SimpleBlockchainService {
	return &SimpleBlockchainService{
		blockchainRepo: blockchainRepo,
		lastBlockHash:  "0", // Genesis block hash
		blockIndex:     0,
	}
}

// RecordTransaction registra uma transação no blockchain
func (s *SimpleBlockchainService) RecordTransaction(ctx context.Context, tx *entities.Transaction) (string, error) {
	// Criar hash da transação
	txHash := s.calculateTransactionHash(tx)

	// Incrementar índice do bloco
	s.blockIndex++

	// Criar registro de blockchain
	record := entities.NewBlockRecord(
		tx.ID,
		s.blockIndex,
		txHash,
		s.lastBlockHash,
		map[string]interface{}{
			"id":          tx.ID,
			"source":      tx.SourceAccount,
			"destination": tx.DestinationAccount,
			"amount":      tx.Amount,
			"status":      tx.Status,
			"recorded_at": time.Now(),
		},
	)

	// Persistir registro
	if err := s.blockchainRepo.SaveRecord(ctx, record); err != nil {
		return "", fmt.Errorf("failed to save blockchain record: %w", err)
	}

	// Atualizar hash anterior
	s.lastBlockHash = txHash

	return txHash, nil
}

// GetTransactionProof obtém o comprovante de uma transação na blockchain
func (s *SimpleBlockchainService) GetTransactionProof(ctx context.Context, txID string) (*entities.BlockRecord, error) {
	return s.blockchainRepo.FindRecordByTransactionID(ctx, txID)
}

// VerifyTransaction verifica se uma transação está na blockchain
func (s *SimpleBlockchainService) VerifyTransaction(ctx context.Context, txID, hash string) (bool, error) {
	record, err := s.blockchainRepo.FindRecordByTransactionID(ctx, txID)
	if err != nil {
		return false, err
	}

	if record == nil {
		return false, nil
	}

	// Verificar integridade do hash
	return record.BlockHash == hash, nil
}

// calculateTransactionHash calcula o hash de uma transação
func (s *SimpleBlockchainService) calculateTransactionHash(tx *entities.Transaction) string {
	data := fmt.Sprintf(
		"%s|%s|%s|%f|%s|%d|%s",
		tx.ID,
		tx.SourceAccount,
		tx.DestinationAccount,
		tx.Amount,
		tx.Currency,
		tx.CreatedAt.Unix(),
		s.lastBlockHash,
	)

	hash := sha256.Sum256([]byte(data))
	return fmt.Sprintf("%x", hash)
}
