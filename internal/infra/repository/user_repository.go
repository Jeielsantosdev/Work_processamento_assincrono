package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Jeielsantosdev/Work_processamento_assincrono/internal/domain/entities"
)

// UserRepositorySQL implementa o repositório de usuários com SQL
type UserRepositorySQL struct {
	db *sql.DB
}

// NewUserRepositorySQL cria uma nova instância do repositório
func NewUserRepositorySQL(db *sql.DB) *UserRepositorySQL {
	return &UserRepositorySQL{db: db}
}

// Save salva um novo usuário
func (r *UserRepositorySQL) Save(ctx context.Context, user *entities.User) error {
	query := `
		INSERT INTO users (id, email, username, password, role, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	_, err := r.db.ExecContext(ctx, query,
		user.ID, user.Email, user.Username, user.Password,
		user.Role, user.IsActive, user.CreatedAt, user.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to save user: %w", err)
	}

	return nil
}

// FindByID busca um usuário por ID
func (r *UserRepositorySQL) FindByID(ctx context.Context, id string) (*entities.User, error) {
	query := `
		SELECT id, email, username, password, role, is_active, created_at, updated_at, last_login
		FROM users WHERE id = $1
	`

	user := &entities.User{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID, &user.Email, &user.Username, &user.Password,
		&user.Role, &user.IsActive, &user.CreatedAt, &user.UpdatedAt, &user.LastLogin,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, entities.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	return user, nil
}

// FindByEmail busca um usuário por email
func (r *UserRepositorySQL) FindByEmail(ctx context.Context, email string) (*entities.User, error) {
	query := `
		SELECT id, email, username, password, role, is_active, created_at, updated_at, last_login
		FROM users WHERE email = $1
	`

	user := &entities.User{}
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID, &user.Email, &user.Username, &user.Password,
		&user.Role, &user.IsActive, &user.CreatedAt, &user.UpdatedAt, &user.LastLogin,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, entities.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	return user, nil
}

// Update atualiza um usuário existente
func (r *UserRepositorySQL) Update(ctx context.Context, user *entities.User) error {
	query := `
		UPDATE users 
		SET email = $1, username = $2, password = $3, role = $4, 
		    is_active = $5, updated_at = $6, last_login = $7
		WHERE id = $8
	`

	result, err := r.db.ExecContext(ctx, query,
		user.Email, user.Username, user.Password, user.Role,
		user.IsActive, user.UpdatedAt, user.LastLogin, user.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return entities.ErrUserNotFound
	}

	return nil
}

// Delete deleta um usuário
func (r *UserRepositorySQL) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM users WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return entities.ErrUserNotFound
	}

	return nil
}
