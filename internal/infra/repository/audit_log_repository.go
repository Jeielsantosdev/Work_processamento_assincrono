package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Jeielsantosdev/Work_processamento_assincrono/internal/domain/entities"
)

// AuditLogRepositorySQL implementa o repositório de logs de auditoria com SQL
type AuditLogRepositorySQL struct {
	db *sql.DB
}

// NewAuditLogRepositorySQL cria uma nova instância do repositório
func NewAuditLogRepositorySQL(db *sql.DB) *AuditLogRepositorySQL {
	return &AuditLogRepositorySQL{db: db}
}

// Save salva um novo log de auditoria
func (r *AuditLogRepositorySQL) Save(ctx context.Context, log *entities.AuditLog) error {
	query := `
		INSERT INTO audit_logs (
			id, transaction_id, action, actor_id, actor_role, 
			ip_address, timestamp, status, error_message
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`

	_, err := r.db.ExecContext(ctx, query,
		log.ID, log.TransactionID, log.Action, log.ActorID, log.ActorRole,
		log.IPAddress, log.Timestamp, log.Status, log.ErrorMessage,
	)

	if err != nil {
		return fmt.Errorf("failed to save audit log: %w", err)
	}

	return nil
}

// FindByTransactionID busca logs por ID de transação
func (r *AuditLogRepositorySQL) FindByTransactionID(ctx context.Context, txID string) ([]*entities.AuditLog, error) {
	query := `
		SELECT id, transaction_id, action, actor_id, actor_role, 
		       ip_address, timestamp, status, error_message
		FROM audit_logs WHERE transaction_id = $1 ORDER BY timestamp DESC
	`

	rows, err := r.db.QueryContext(ctx, query, txID)
	if err != nil {
		return nil, fmt.Errorf("failed to query audit logs: %w", err)
	}
	defer rows.Close()

	var logs []*entities.AuditLog
	for rows.Next() {
		log := &entities.AuditLog{}
		err := rows.Scan(
			&log.ID, &log.TransactionID, &log.Action, &log.ActorID, &log.ActorRole,
			&log.IPAddress, &log.Timestamp, &log.Status, &log.ErrorMessage,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan audit log: %w", err)
		}
		logs = append(logs, log)
	}

	return logs, nil
}

// FindByActorID busca logs por ID do ator
func (r *AuditLogRepositorySQL) FindByActorID(ctx context.Context, actorID string) ([]*entities.AuditLog, error) {
	query := `
		SELECT id, transaction_id, action, actor_id, actor_role, 
		       ip_address, timestamp, status, error_message
		FROM audit_logs WHERE actor_id = $1 ORDER BY timestamp DESC
	`

	rows, err := r.db.QueryContext(ctx, query, actorID)
	if err != nil {
		return nil, fmt.Errorf("failed to query audit logs: %w", err)
	}
	defer rows.Close()

	var logs []*entities.AuditLog
	for rows.Next() {
		log := &entities.AuditLog{}
		err := rows.Scan(
			&log.ID, &log.TransactionID, &log.Action, &log.ActorID, &log.ActorRole,
			&log.IPAddress, &log.Timestamp, &log.Status, &log.ErrorMessage,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan audit log: %w", err)
		}
		logs = append(logs, log)
	}

	return logs, nil
}

// ListAll lista todos os logs de auditoria
func (r *AuditLogRepositorySQL) ListAll(ctx context.Context, limit, offset int) ([]*entities.AuditLog, error) {
	query := `
		SELECT id, transaction_id, action, actor_id, actor_role, 
		       ip_address, timestamp, status, error_message
		FROM audit_logs ORDER BY timestamp DESC LIMIT $1 OFFSET $2
	`

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to query audit logs: %w", err)
	}
	defer rows.Close()

	var logs []*entities.AuditLog
	for rows.Next() {
		log := &entities.AuditLog{}
		err := rows.Scan(
			&log.ID, &log.TransactionID, &log.Action, &log.ActorID, &log.ActorRole,
			&log.IPAddress, &log.Timestamp, &log.Status, &log.ErrorMessage,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan audit log: %w", err)
		}
		logs = append(logs, log)
	}

	return logs, nil
}
