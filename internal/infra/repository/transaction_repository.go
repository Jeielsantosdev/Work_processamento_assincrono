package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Jeielsantosdev/Work_processamento_assincrono/internal/domain/entities"
)

// TransactionRepositorySQL implementa o repositório de transações com SQL
type TransactionRepositorySQL struct {
	db *sql.DB
}

// NewTransactionRepositorySQL cria uma nova instância do repositório
func NewTransactionRepositorySQL(db *sql.DB) *TransactionRepositorySQL {
	return &TransactionRepositorySQL{db: db}
}

// Save salva uma nova transação
func (r *TransactionRepositorySQL) Save(ctx context.Context, transaction *entities.Transaction) error {
	query := `
		INSERT INTO transactions (
			id, client_id, source_account, destination_account, amount, 
			currency, description, status, retry_count, max_retries, 
			created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	`

	_, err := r.db.ExecContext(ctx, query,
		transaction.ID,
		transaction.ClientID,
		transaction.SourceAccount,
		transaction.DestinationAccount,
		transaction.Amount,
		transaction.Currency,
		transaction.Description,
		transaction.Status,
		transaction.RetryCount,
		transaction.MaxRetries,
		transaction.CreatedAt,
		transaction.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to save transaction: %w", err)
	}

	return nil
}

// FindByID busca uma transação por ID
func (r *TransactionRepositorySQL) FindByID(ctx context.Context, id string) (*entities.Transaction, error) {
	query := `
		SELECT id, client_id, source_account, destination_account, amount,
		       currency, description, status, reason, retry_count, 
		       max_retries, blockchain_hash, created_at, updated_at, processed_at
		FROM transactions WHERE id = $1
	`

	tx := &entities.Transaction{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&tx.ID, &tx.ClientID, &tx.SourceAccount, &tx.DestinationAccount,
		&tx.Amount, &tx.Currency, &tx.Description, &tx.Status, &tx.Reason,
		&tx.RetryCount, &tx.MaxRetries, &tx.BlockchainHash,
		&tx.CreatedAt, &tx.UpdatedAt, &tx.ProcessedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, entities.ErrTransactionNotFound
		}
		return nil, fmt.Errorf("failed to find transaction: %w", err)
	}

	return tx, nil
}

// FindByClientID busca transações de um cliente
func (r *TransactionRepositorySQL) FindByClientID(ctx context.Context, clientID string) ([]*entities.Transaction, error) {
	query := `
		SELECT id, client_id, source_account, destination_account, amount,
		       currency, description, status, reason, retry_count, 
		       max_retries, blockchain_hash, created_at, updated_at, processed_at
		FROM transactions WHERE client_id = $1 ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, clientID)
	if err != nil {
		return nil, fmt.Errorf("failed to query transactions: %w", err)
	}
	defer rows.Close()

	var transactions []*entities.Transaction
	for rows.Next() {
		tx := &entities.Transaction{}
		err := rows.Scan(
			&tx.ID, &tx.ClientID, &tx.SourceAccount, &tx.DestinationAccount,
			&tx.Amount, &tx.Currency, &tx.Description, &tx.Status, &tx.Reason,
			&tx.RetryCount, &tx.MaxRetries, &tx.BlockchainHash,
			&tx.CreatedAt, &tx.UpdatedAt, &tx.ProcessedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan transaction: %w", err)
		}
		transactions = append(transactions, tx)
	}

	return transactions, nil
}

// Update atualiza uma transação existente
func (r *TransactionRepositorySQL) Update(ctx context.Context, transaction *entities.Transaction) error {
	query := `
		UPDATE transactions 
		SET status = $1, reason = $2, retry_count = $3, 
		    blockchain_hash = $4, updated_at = $5, processed_at = $6
		WHERE id = $7
	`

	result, err := r.db.ExecContext(ctx, query,
		transaction.Status, transaction.Reason, transaction.RetryCount,
		transaction.BlockchainHash, transaction.UpdatedAt, transaction.ProcessedAt,
		transaction.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update transaction: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return entities.ErrTransactionNotFound
	}

	return nil
}

// ListAll lista todas as transações
func (r *TransactionRepositorySQL) ListAll(ctx context.Context, limit, offset int) ([]*entities.Transaction, error) {
	query := `
		SELECT id, client_id, source_account, destination_account, amount,
		       currency, description, status, reason, retry_count, 
		       max_retries, blockchain_hash, created_at, updated_at, processed_at
		FROM transactions ORDER BY created_at DESC LIMIT $1 OFFSET $2
	`

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to query transactions: %w", err)
	}
	defer rows.Close()

	var transactions []*entities.Transaction
	for rows.Next() {
		tx := &entities.Transaction{}
		err := rows.Scan(
			&tx.ID, &tx.ClientID, &tx.SourceAccount, &tx.DestinationAccount,
			&tx.Amount, &tx.Currency, &tx.Description, &tx.Status, &tx.Reason,
			&tx.RetryCount, &tx.MaxRetries, &tx.BlockchainHash,
			&tx.CreatedAt, &tx.UpdatedAt, &tx.ProcessedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan transaction: %w", err)
		}
		transactions = append(transactions, tx)
	}

	return transactions, nil
}
