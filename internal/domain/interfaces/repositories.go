package interfaces

import (
	"context"

	"github.com/Jeielsantosdev/Work_processamento_assincrono/internal/domain/entities"
)

// TransactionRepository define o contrato para persistência de transações
type TransactionRepository interface {
	Save(ctx context.Context, transaction *entities.Transaction) error
	FindByID(ctx context.Context, id string) (*entities.Transaction, error)
	FindByClientID(ctx context.Context, clientID string) ([]*entities.Transaction, error)
	Update(ctx context.Context, transaction *entities.Transaction) error
	ListAll(ctx context.Context, limit, offset int) ([]*entities.Transaction, error)
}

// UserRepository define o contrato para persistência de usuários
type UserRepository interface {
	Save(ctx context.Context, user *entities.User) error
	FindByID(ctx context.Context, id string) (*entities.User, error)
	FindByEmail(ctx context.Context, email string) (*entities.User, error)
	Update(ctx context.Context, user *entities.User) error
	Delete(ctx context.Context, id string) error
}

// NotificationRepository define o contrato para persistência de notificações
type NotificationRepository interface {
	Save(ctx context.Context, notification *entities.Notification) error
	FindByID(ctx context.Context, id string) (*entities.Notification, error)
	FindByClientID(ctx context.Context, clientID string) ([]*entities.Notification, error)
	Update(ctx context.Context, notification *entities.Notification) error
}

// BlockchainRepository define o contrato para persistência de registros blockchain
type BlockchainRepository interface {
	SaveRecord(ctx context.Context, record *entities.BlockRecord) error
	FindRecordByTransactionID(ctx context.Context, txID string) (*entities.BlockRecord, error)
	GetLatestBlock(ctx context.Context) (*entities.Block, error)
}

// AuditLogRepository define o contrato para persistência de logs de auditoria
type AuditLogRepository interface {
	Save(ctx context.Context, log *entities.AuditLog) error
	FindByTransactionID(ctx context.Context, txID string) ([]*entities.AuditLog, error)
	FindByActorID(ctx context.Context, actorID string) ([]*entities.AuditLog, error)
	ListAll(ctx context.Context, limit, offset int) ([]*entities.AuditLog, error)
}
