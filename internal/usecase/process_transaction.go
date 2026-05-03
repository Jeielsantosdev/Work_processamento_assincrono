package usecase

import (
	"context"
	"fmt"

	"github.com/Jeielsantosdev/Work_processamento_assincrono/internal/domain/entities"
	"github.com/Jeielsantosdev/Work_processamento_assincrono/internal/domain/interfaces"
)

// ProcessTransactionUseCase implementa o caso de uso de processar transação
type ProcessTransactionUseCase struct {
	transactionRepo interfaces.TransactionRepository
	blockchainSvc   interfaces.BlockchainService
	notificationSvc interfaces.NotificationService
	auditLogRepo    interfaces.AuditLogRepository
}

// NewProcessTransactionUseCase cria uma nova instância do use case
func NewProcessTransactionUseCase(
	transactionRepo interfaces.TransactionRepository,
	blockchainSvc interfaces.BlockchainService,
	notificationSvc interfaces.NotificationService,
	auditLogRepo interfaces.AuditLogRepository,
) *ProcessTransactionUseCase {
	return &ProcessTransactionUseCase{
		transactionRepo: transactionRepo,
		blockchainSvc:   blockchainSvc,
		notificationSvc: notificationSvc,
		auditLogRepo:    auditLogRepo,
	}
}

// Execute executa o processamento da transação
func (uc *ProcessTransactionUseCase) Execute(ctx context.Context, transactionID string) error {
	// Recuperar transação
	tx, err := uc.transactionRepo.FindByID(ctx, transactionID)
	if err != nil {
		return fmt.Errorf("failed to find transaction: %w", err)
	}

	// Validar se pode processar
	if tx.Status != entities.StatusPending {
		return fmt.Errorf("transaction status is %s, cannot process", tx.Status)
	}

	// Marcar como em processamento
	tx.MarkAsProcessing()
	if err := uc.transactionRepo.Update(ctx, tx); err != nil {
		return fmt.Errorf("failed to update transaction: %w", err)
	}

	// Validações de negócio (exemplo: verificar saldo)
	if err := uc.validateBusinessRules(ctx, tx); err != nil {
		tx.MarkAsFailed(err.Error())
		uc.transactionRepo.Update(ctx, tx)
		uc.notifyTransactionFailure(ctx, tx)
		return nil // Não retorna erro, pois é falha de negócio
	}

	// Registrar na blockchain
	blockchainHash, err := uc.blockchainSvc.RecordTransaction(ctx, tx)
	if err != nil {
		tx.MarkAsFailed("Failed to record in blockchain")
		uc.transactionRepo.Update(ctx, tx)
		uc.notifyTransactionFailure(ctx, tx)
		return nil
	}

	// Marcar como completa
	tx.MarkAsCompleted(blockchainHash)
	if err := uc.transactionRepo.Update(ctx, tx); err != nil {
		return fmt.Errorf("failed to finalize transaction: %w", err)
	}

	// Notificar cliente
	uc.notifyTransactionSuccess(ctx, tx)

	// Log de auditoria
	auditLog := entities.NewAuditLog(
		tx.ID,
		"TRANSACTION_PROCESSED",
		"SYSTEM",
		entities.RoleAdministrator,
		map[string]interface{}{
			"blockchain_hash": blockchainHash,
			"amount":          tx.Amount,
		},
		"INTERNAL",
	)
	uc.auditLogRepo.Save(ctx, auditLog)

	return nil
}

func (uc *ProcessTransactionUseCase) validateBusinessRules(ctx context.Context, tx *entities.Transaction) error {
	// Validar saldo (integrar com serviço externo de conta corrente)
	// Validar limites
	// Validar horários
	// Etc...
	return nil
}

func (uc *ProcessTransactionUseCase) notifyTransactionSuccess(ctx context.Context, tx *entities.Transaction) {
	notification := entities.NewNotification(
		tx.ClientID,
		entities.NotificationTypeTransactionStatus,
		"Transação Concluída",
		fmt.Sprintf("Sua transação de R$ %.2f foi processada com sucesso", tx.Amount),
		map[string]interface{}{
			"transaction_id": tx.ID,
			"status":         tx.Status,
			"amount":         tx.Amount,
		},
		"email",
	)
	uc.notificationSvc.SendNotification(ctx, notification)
}

func (uc *ProcessTransactionUseCase) notifyTransactionFailure(ctx context.Context, tx *entities.Transaction) {
	notification := entities.NewNotification(
		tx.ClientID,
		entities.NotificationTypeTransactionStatus,
		"Transação Falhada",
		fmt.Sprintf("Sua transação não pôde ser processada: %s", tx.Reason),
		map[string]interface{}{
			"transaction_id": tx.ID,
			"status":         tx.Status,
			"reason":         tx.Reason,
		},
		"email",
	)
	uc.notificationSvc.SendNotification(ctx, notification)
}
