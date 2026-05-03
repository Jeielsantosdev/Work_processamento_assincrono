package usecase

import (
	"context"
	"fmt"

	"github.com/Jeielsantosdev/Work_processamento_assincrono/internal/domain/entities"
	"github.com/Jeielsantosdev/Work_processamento_assincrono/internal/domain/interfaces"
)

// CreateTransactionUseCase implementa o caso de uso de criar transação
type CreateTransactionUseCase struct {
	transactionRepo interfaces.TransactionRepository
	queueService    interfaces.QueueService
	auditLogRepo    interfaces.AuditLogRepository
}

// NewCreateTransactionUseCase cria uma nova instância do use case
func NewCreateTransactionUseCase(
	transactionRepo interfaces.TransactionRepository,
	queueService interfaces.QueueService,
	auditLogRepo interfaces.AuditLogRepository,
) *CreateTransactionUseCase {
	return &CreateTransactionUseCase{
		transactionRepo: transactionRepo,
		queueService:    queueService,
		auditLogRepo:    auditLogRepo,
	}
}

// Execute executa o caso de uso de criação de transação
func (uc *CreateTransactionUseCase) Execute(ctx context.Context, input *CreateTransactionInput) (*CreateTransactionOutput, error) {
	// Validar entrada
	if err := input.Validate(); err != nil {
		return nil, fmt.Errorf("invalid input: %w", err)
	}

	// Criar nova transação
	transaction, err := entities.NewTransaction(
		input.ClientID,
		input.SourceAccount,
		input.DestinationAccount,
		input.Amount,
		input.Currency,
		input.Description,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create transaction: %w", err)
	}

	// Persitir transação
	if err := uc.transactionRepo.Save(ctx, transaction); err != nil {
		return nil, fmt.Errorf("failed to save transaction: %w", err)
	}

	// Publicar na fila para processamento assíncrono
	if err := uc.queueService.PublishTransaction(ctx, transaction); err != nil {
		// Marcar como falha
		transaction.MarkAsFailed("Failed to enqueue transaction")
		uc.transactionRepo.Update(ctx, transaction)
		return nil, fmt.Errorf("failed to publish transaction to queue: %w", err)
	}

	// Registrar log de auditoria
	auditLog := entities.NewAuditLog(
		transaction.ID,
		"TRANSACTION_CREATED",
		input.ClientID,
		entities.RoleClient,
		map[string]interface{}{
			"amount":      input.Amount,
			"source":      input.SourceAccount,
			"destination": input.DestinationAccount,
		},
		input.IPAddress,
	)
	uc.auditLogRepo.Save(ctx, auditLog)

	return &CreateTransactionOutput{
		TransactionID: transaction.ID,
		Status:        string(transaction.Status),
		CreatedAt:     transaction.CreatedAt.String(),
	}, nil
}

// CreateTransactionInput representa a entrada do use case
type CreateTransactionInput struct {
	ClientID           string
	SourceAccount      string
	DestinationAccount string
	Amount             float64
	Currency           string
	Description        string
	IPAddress          string
}

// Validate valida a entrada
func (i *CreateTransactionInput) Validate() error {
	if i.ClientID == "" {
		return entities.ErrInvalidClientID
	}
	if i.SourceAccount == "" || i.DestinationAccount == "" {
		return entities.ErrInvalidAccount
	}
	if i.Amount <= 0 {
		return entities.ErrInvalidAmount
	}
	return nil
}

// CreateTransactionOutput representa a saída do use case
type CreateTransactionOutput struct {
	TransactionID string
	Status        string
	CreatedAt     string
}
