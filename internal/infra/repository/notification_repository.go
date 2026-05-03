package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Jeielsantosdev/Work_processamento_assincrono/internal/domain/entities"
)

// NotificationRepositorySQL implementa o repositório de notificações com SQL
type NotificationRepositorySQL struct {
	db *sql.DB
}

// NewNotificationRepositorySQL cria uma nova instância do repositório
func NewNotificationRepositorySQL(db *sql.DB) *NotificationRepositorySQL {
	return &NotificationRepositorySQL{db: db}
}

// Save salva uma nova notificação
func (r *NotificationRepositorySQL) Save(ctx context.Context, notification *entities.Notification) error {
	query := `
		INSERT INTO notifications (
			id, client_id, type, title, message, status, channel, 
			created_at, updated_at, sent_at, failure_reason
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := r.db.ExecContext(ctx, query,
		notification.ID, notification.ClientID, notification.Type, notification.Title,
		notification.Message, notification.Status, notification.Channel,
		notification.CreatedAt, notification.UpdatedAt, notification.SentAt, notification.FailureReason,
	)

	if err != nil {
		return fmt.Errorf("failed to save notification: %w", err)
	}

	return nil
}

// FindByID busca uma notificação por ID
func (r *NotificationRepositorySQL) FindByID(ctx context.Context, id string) (*entities.Notification, error) {
	query := `
		SELECT id, client_id, type, title, message, status, channel, 
		       created_at, updated_at, sent_at, failure_reason
		FROM notifications WHERE id = ?
	`

	notification := &entities.Notification{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&notification.ID, &notification.ClientID, &notification.Type, &notification.Title,
		&notification.Message, &notification.Status, &notification.Channel,
		&notification.CreatedAt, &notification.UpdatedAt, &notification.SentAt, &notification.FailureReason,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find notification: %w", err)
	}

	return notification, nil
}

// FindByClientID busca notificações de um cliente
func (r *NotificationRepositorySQL) FindByClientID(ctx context.Context, clientID string) ([]*entities.Notification, error) {
	query := `
		SELECT id, client_id, type, title, message, status, channel, 
		       created_at, updated_at, sent_at, failure_reason
		FROM notifications WHERE client_id = ? ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, clientID)
	if err != nil {
		return nil, fmt.Errorf("failed to query notifications: %w", err)
	}
	defer rows.Close()

	var notifications []*entities.Notification
	for rows.Next() {
		notification := &entities.Notification{}
		err := rows.Scan(
			&notification.ID, &notification.ClientID, &notification.Type, &notification.Title,
			&notification.Message, &notification.Status, &notification.Channel,
			&notification.CreatedAt, &notification.UpdatedAt, &notification.SentAt, &notification.FailureReason,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan notification: %w", err)
		}
		notifications = append(notifications, notification)
	}

	return notifications, nil
}

// Update atualiza uma notificação existente
func (r *NotificationRepositorySQL) Update(ctx context.Context, notification *entities.Notification) error {
	query := `
		UPDATE notifications 
		SET status = ?, sent_at = ?, failure_reason = ?, updated_at = ?
		WHERE id = ?
	`

	result, err := r.db.ExecContext(ctx, query,
		notification.Status, notification.SentAt, notification.FailureReason, notification.UpdatedAt,
		notification.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update notification: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("notification not found")
	}

	return nil
}
