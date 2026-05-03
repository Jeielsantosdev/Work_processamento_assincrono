package interfaces

import (
	"context"

	"github.com/Jeielsantosdev/Work_processamento_assincrono/internal/domain/entities"
)

// BlockchainService define o contrato para operações de blockchain
type BlockchainService interface {
	RecordTransaction(ctx context.Context, tx *entities.Transaction) (string, error) // Retorna hash
	GetTransactionProof(ctx context.Context, txID string) (*entities.BlockRecord, error)
	VerifyTransaction(ctx context.Context, txID, hash string) (bool, error)
}

// NotificationService define o contrato para envio de notificações
type NotificationService interface {
	SendNotification(ctx context.Context, notification *entities.Notification) error
	BroadcastTransactionStatus(ctx context.Context, tx *entities.Transaction) error
}

// QueueService define o contrato para operações de fila
type QueueService interface {
	PublishTransaction(ctx context.Context, tx *entities.Transaction) error
	ConsumeTransactions(ctx context.Context) (<-chan *entities.Transaction, error)
	AcknowledgeMessage(ctx context.Context, messageID string) error
}

// PasswordService define o contrato para operações de senha
type PasswordService interface {
	HashPassword(password string) (string, error)
	VerifyPassword(hashedPassword, plainPassword string) bool
}

// AuthenticationService define o contrato para autenticação
type AuthenticationService interface {
	GenerateToken(user *entities.User) (string, error)
	ValidateToken(token string) (*TokenClaims, error)
}

// TokenClaims representa as claims de um JWT
type TokenClaims struct {
	UserID    string
	Email     string
	Role      entities.UserRole
	IssuedAt  int64
	ExpiresAt int64
}
