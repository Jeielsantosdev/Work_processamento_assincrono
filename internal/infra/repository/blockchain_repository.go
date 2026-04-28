package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Jeielsantosdev/Work_processamento_assincrono/internal/domain/entities"
)

// BlockchainRepositorySQL implementa o repositório de blockchain com SQL
type BlockchainRepositorySQL struct {
	db *sql.DB
}

// NewBlockchainRepositorySQL cria uma nova instância do repositório
func NewBlockchainRepositorySQL(db *sql.DB) *BlockchainRepositorySQL {
	return &BlockchainRepositorySQL{db: db}
}

// SaveRecord salva um registro de blockchain
func (r *BlockchainRepositorySQL) SaveRecord(ctx context.Context, record *entities.BlockRecord) error {
	query := `
		INSERT INTO blockchain_records (
			id, transaction_id, block_index, block_hash, previous_hash, 
			timestamp, status, transaction_data
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	_, err := r.db.ExecContext(ctx, query,
		record.ID, record.TransactionID, record.BlockIndex, record.BlockHash,
		record.PreviousHash, record.Timestamp, record.Status,
		// TODO: Serializar transaction_data como JSON
		"{}",
	)

	if err != nil {
		return fmt.Errorf("failed to save blockchain record: %w", err)
	}

	return nil
}

// FindRecordByTransactionID busca um registro por ID de transação
func (r *BlockchainRepositorySQL) FindRecordByTransactionID(ctx context.Context, txID string) (*entities.BlockRecord, error) {
	query := `
		SELECT id, transaction_id, block_index, block_hash, previous_hash, 
		       timestamp, status, transaction_data
		FROM blockchain_records WHERE transaction_id = $1 LIMIT 1
	`

	record := &entities.BlockRecord{}
	var jsonData string

	err := r.db.QueryRowContext(ctx, query, txID).Scan(
		&record.ID, &record.TransactionID, &record.BlockIndex, &record.BlockHash,
		&record.PreviousHash, &record.Timestamp, &record.Status, &jsonData,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find blockchain record: %w", err)
	}

	// TODO: Desserializar JSON em record.TransactionData

	return record, nil
}

// GetLatestBlock retorna o último bloco
func (r *BlockchainRepositorySQL) GetLatestBlock(ctx context.Context) (*entities.Block, error) {
	// TODO: Implementar lógica para retornar o último bloco
	return nil, nil
}
