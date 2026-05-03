package entities

import (
	"time"

	"github.com/google/uuid"
)

// TransactionStatus representa o status de uma transação
type TransactionStatus string

const (
	StatusPending    TransactionStatus = "PENDING"
	StatusProcessing TransactionStatus = "PROCESSING"
	StatusCompleted  TransactionStatus = "COMPLETED"
	StatusFailed     TransactionStatus = "FAILED"
	StatusRejected   TransactionStatus = "REJECTED"
)

// Transaction representa uma transação financeira no sistema
type Transaction struct {
	ID                 string
	ClientID           string
	SourceAccount      string
	DestinationAccount string
	Amount             float64
	Currency           string
	Description        string
	Status             TransactionStatus
	Reason             string // Motivo da falha ou rejeição
	RetryCount         int
	MaxRetries         int
	BlockchainHash     string
	CreatedAt          time.Time
	UpdatedAt          time.Time
	ProcessedAt        *time.Time
}

// NewTransaction cria uma nova transação com validação básica
func NewTransaction(clientID, sourceAccount, destAccount string, amount float64, currency, description string) (*Transaction, error) {
	if err := validateTransaction(clientID, sourceAccount, destAccount, amount); err != nil {
		return nil, err
	}

	return &Transaction{
		ID:                 uuid.New().String(),
		ClientID:           clientID,
		SourceAccount:      sourceAccount,
		DestinationAccount: destAccount,
		Amount:             amount,
		Currency:           currency,
		Description:        description,
		Status:             StatusPending,
		RetryCount:         0,
		MaxRetries:         3,
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
	}, nil
}

// IsValid verifica se a transação está em estado válido
func (t *Transaction) IsValid() bool {
	return t.ID != "" && t.ClientID != "" && t.Amount > 0 && t.Status != ""
}

// CanRetry verifica se a transação pode ser reprocessada
func (t *Transaction) CanRetry() bool {
	return t.Status == StatusFailed && t.RetryCount < t.MaxRetries
}

// IncrementRetry incrementa o contador de tentativas
func (t *Transaction) IncrementRetry() {
	t.RetryCount++
	t.UpdatedAt = time.Now()
}

// MarkAsProcessing marca a transação como em processamento
func (t *Transaction) MarkAsProcessing() {
	t.Status = StatusProcessing
	t.UpdatedAt = time.Now()
}

// MarkAsCompleted marca a transação como completada
func (t *Transaction) MarkAsCompleted(blockchainHash string) {
	t.Status = StatusCompleted
	t.BlockchainHash = blockchainHash
	now := time.Now()
	t.ProcessedAt = &now
	t.UpdatedAt = now
}

// MarkAsFailed marca a transação como falha
func (t *Transaction) MarkAsFailed(reason string) {
	t.Status = StatusFailed
	t.Reason = reason
	t.UpdatedAt = time.Now()
}

// MarkAsRejected marca a transação como rejeitada
func (t *Transaction) MarkAsRejected(reason string) {
	t.Status = StatusRejected
	t.Reason = reason
	t.UpdatedAt = time.Now()
}

func validateTransaction(clientID, source, dest string, amount float64) error {
	if clientID == "" {
		return ErrInvalidClientID
	}
	if source == "" || dest == "" {
		return ErrInvalidAccount
	}
	if amount <= 0 {
		return ErrInvalidAmount
	}
	if source == dest {
		return ErrSameAccount
	}
	return nil
}
