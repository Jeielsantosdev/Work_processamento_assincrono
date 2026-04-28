package entities

import (
	"time"

	"github.com/google/uuid"
)

// UserRole representa o papel do usuário no sistema
type UserRole string

const (
	RoleClient        UserRole = "CLIENT"
	RoleAuditor       UserRole = "AUDITOR"
	RoleAdministrator UserRole = "ADMINISTRATOR"
)

// User representa um usuário do sistema
type User struct {
	ID        string
	Email     string
	Username  string
	Password  string // Hash da senha
	Role      UserRole
	IsActive  bool
	CreatedAt time.Time
	UpdatedAt time.Time
	LastLogin *time.Time
}

// NewUser cria um novo usuário
func NewUser(email, username, passwordHash string, role UserRole) *User {
	return &User{
		ID:        uuid.New().String(),
		Email:     email,
		Username:  username,
		Password:  passwordHash,
		Role:      role,
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// HasPermission verifica se o usuário tem permissão para uma ação
func (u *User) HasPermission(action string) bool {
	if !u.IsActive {
		return false
	}

	permissions := map[UserRole][]string{
		RoleClient: {
			"send_transaction",
			"view_own_transactions",
		},
		RoleAuditor: {
			"view_all_transactions",
			"generate_reports",
			"view_blockchain",
		},
		RoleAdministrator: {
			"manage_users",
			"manage_system",
			"view_all_transactions",
			"manage_workers",
			"manage_queue",
		},
	}

	allowedActions, exists := permissions[u.Role]
	if !exists {
		return false
	}

	for _, allowed := range allowedActions {
		if allowed == action {
			return true
		}
	}
	return false
}

// UpdateLastLogin atualiza o timestamp do último login
func (u *User) UpdateLastLogin() {
	now := time.Now()
	u.LastLogin = &now
	u.UpdatedAt = now
}

// Deactivate desativa o usuário
func (u *User) Deactivate() {
	u.IsActive = false
	u.UpdatedAt = time.Now()
}
