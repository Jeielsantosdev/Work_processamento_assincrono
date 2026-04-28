package notification

import (
	"context"
	"fmt"

	"github.com/Jeielsantosdev/Work_processamento_assincrono/internal/domain/entities"
	"github.com/Jeielsantosdev/Work_processamento_assincrono/internal/domain/interfaces"
)

// ConsoleNotificationService implementa o serviço de notificação por console (para desenvolvimento)
type ConsoleNotificationService struct {
	notificationRepo interfaces.NotificationRepository
}

// NewConsoleNotificationService cria uma nova instância do serviço
func NewConsoleNotificationService(notificationRepo interfaces.NotificationRepository) *ConsoleNotificationService {
	return &ConsoleNotificationService{
		notificationRepo: notificationRepo,
	}
}

// SendNotification envia uma notificação
func (s *ConsoleNotificationService) SendNotification(ctx context.Context, notification *entities.Notification) error {
	// Simular envio (em produção seria integrado com Sendgrid, AWS SES, etc)
	fmt.Printf("[NOTIFICATION] To: %s | Title: %s | Message: %s | Channel: %s\n",
		notification.ClientID,
		notification.Title,
		notification.Message,
		notification.Channel,
	)

	// Marcar como enviado
	notification.MarkAsSent()

	// Persistir resultado
	if err := s.notificationRepo.Update(ctx, notification); err != nil {
		notification.MarkAsFailed(err.Error())
		return fmt.Errorf("failed to save notification: %w", err)
	}

	return nil
}

// BroadcastTransactionStatus envia o status de uma transação
func (s *ConsoleNotificationService) BroadcastTransactionStatus(ctx context.Context, tx *entities.Transaction) error {
	var title string
	var message string

	switch tx.Status {
	case entities.StatusCompleted:
		title = "Transação Concluída"
		message = fmt.Sprintf("Sua transação de R$ %.2f foi processada com sucesso", tx.Amount)
	case entities.StatusFailed:
		title = "Transação Falhada"
		message = fmt.Sprintf("Sua transação falhou: %s", tx.Reason)
	case entities.StatusRejected:
		title = "Transação Rejeitada"
		message = fmt.Sprintf("Sua transação foi rejeitada: %s", tx.Reason)
	default:
		return nil
	}

	notification := entities.NewNotification(
		tx.ClientID,
		entities.NotificationTypeTransactionStatus,
		title,
		message,
		map[string]interface{}{
			"transaction_id": tx.ID,
			"amount":         tx.Amount,
			"status":         tx.Status,
		},
		"email",
	)

	return s.SendNotification(ctx, notification)
}
