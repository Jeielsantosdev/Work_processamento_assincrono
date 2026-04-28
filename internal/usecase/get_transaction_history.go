package usecase

import (
	"context"
	"fmt"

	"github.com/Jeielsantosdev/Work_processamento_assincrono/internal/domain/entities"
	"github.com/Jeielsantosdev/Work_processamento_assincrono/internal/domain/interfaces"
)

// GetTransactionHistoryUseCase implementa o caso de uso de obter histórico
type GetTransactionHistoryUseCase struct {
	transactionRepo interfaces.TransactionRepository
	auditLogRepo    interfaces.AuditLogRepository
}

// NewGetTransactionHistoryUseCase cria uma nova instância do use case
func NewGetTransactionHistoryUseCase(
	transactionRepo interfaces.TransactionRepository,
	auditLogRepo interfaces.AuditLogRepository,
) *GetTransactionHistoryUseCase {
	return &GetTransactionHistoryUseCase{
		transactionRepo: transactionRepo,
		auditLogRepo:    auditLogRepo,
	}
}

// ExecuteForClient executa a busca de histórico para um cliente
func (uc *GetTransactionHistoryUseCase) ExecuteForClient(ctx context.Context, clientID string, user *entities.User) (*GetTransactionHistoryOutput, error) {
	// Verificar permissões
	if user.Role == entities.RoleClient && user.ID != clientID {
		return nil, entities.ErrUnauthorized
	}

	// Recuperar transações
	transactions, err := uc.transactionRepo.FindByClientID(ctx, clientID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch transactions: %w", err)
	}

	// Recuperar logs de auditoria
	auditLogs, err := uc.auditLogRepo.FindByActorID(ctx, clientID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch audit logs: %w", err)
	}

	output := &GetTransactionHistoryOutput{
		Transactions: make([]TransactionDTO, 0, len(transactions)),
		AuditLogs:    make([]AuditLogDTO, 0, len(auditLogs)),
	}

	// Mapear transações para DTO
	for _, tx := range transactions {
		output.Transactions = append(output.Transactions, TransactionDTO{
			ID:                 tx.ID,
			SourceAccount:      tx.SourceAccount,
			DestinationAccount: tx.DestinationAccount,
			Amount:             tx.Amount,
			Currency:           tx.Currency,
			Status:             string(tx.Status),
			CreatedAt:          tx.CreatedAt.String(),
			ProcessedAt:        tx.ProcessedAt.String(),
		})
	}

	// Mapear logs para DTO
	for _, log := range auditLogs {
		output.AuditLogs = append(output.AuditLogs, AuditLogDTO{
			ID:        log.ID,
			Action:    log.Action,
			Timestamp: log.Timestamp.String(),
			Status:    log.Status,
		})
	}

	return output, nil
}

// ExecuteForAuditor executa a busca de histórico para um auditor (acesso a todas as transações)
func (uc *GetTransactionHistoryUseCase) ExecuteForAuditor(ctx context.Context, limit, offset int) (*GetAllTransactionsOutput, error) {
	transactions, err := uc.transactionRepo.ListAll(ctx, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch transactions: %w", err)
	}

	output := &GetAllTransactionsOutput{
		Transactions: make([]TransactionDTO, 0, len(transactions)),
	}

	for _, tx := range transactions {
		output.Transactions = append(output.Transactions, TransactionDTO{
			ID:                 tx.ID,
			SourceAccount:      tx.SourceAccount,
			DestinationAccount: tx.DestinationAccount,
			Amount:             tx.Amount,
			Currency:           tx.Currency,
			Status:             string(tx.Status),
			CreatedAt:          tx.CreatedAt.String(),
			ProcessedAt:        tx.ProcessedAt.String(),
		})
	}

	return output, nil
}

// GetTransactionHistoryOutput representa a saída para cliente
type GetTransactionHistoryOutput struct {
	Transactions []TransactionDTO
	AuditLogs    []AuditLogDTO
}

// GetAllTransactionsOutput representa a saída para auditor
type GetAllTransactionsOutput struct {
	Transactions []TransactionDTO
	Total        int
}

// TransactionDTO é o Data Transfer Object para transação
type TransactionDTO struct {
	ID                 string
	SourceAccount      string
	DestinationAccount string
	Amount             float64
	Currency           string
	Status             string
	CreatedAt          string
	ProcessedAt        string
}

// AuditLogDTO é o Data Transfer Object para log de auditoria
type AuditLogDTO struct {
	ID        string
	Action    string
	Timestamp string
	Status    string
}
