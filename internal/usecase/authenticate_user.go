package usecase

import (
	"context"
	"fmt"

	"github.com/Jeielsantosdev/Work_processamento_assincrono/internal/domain/entities"
	"github.com/Jeielsantosdev/Work_processamento_assincrono/internal/domain/interfaces"
)

// AuthenticateUserUseCase implementa o caso de uso de autenticação
type AuthenticateUserUseCase struct {
	userRepo    interfaces.UserRepository
	passwordSvc interfaces.PasswordService
	authSvc     interfaces.AuthenticationService
}

// NewAuthenticateUserUseCase cria uma nova instância do use case
func NewAuthenticateUserUseCase(
	userRepo interfaces.UserRepository,
	passwordSvc interfaces.PasswordService,
	authSvc interfaces.AuthenticationService,
) *AuthenticateUserUseCase {
	return &AuthenticateUserUseCase{
		userRepo:    userRepo,
		passwordSvc: passwordSvc,
		authSvc:     authSvc,
	}
}

// Execute executa a autenticação do usuário
func (uc *AuthenticateUserUseCase) Execute(ctx context.Context, email, password string) (*AuthenticateOutput, error) {
	// Validar entrada
	if email == "" || password == "" {
		return nil, entities.ErrInvalidCredentials
	}

	// Buscar usuário pelo email
	user, err := uc.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, entities.ErrUserNotFound
	}

	// Verificar se usuário está ativo
	if !user.IsActive {
		return nil, entities.ErrUnauthorized
	}

	// Verificar senha
	if !uc.passwordSvc.VerifyPassword(user.Password, password) {
		return nil, entities.ErrInvalidCredentials
	}

	// Atualizar último login
	user.UpdateLastLogin()
	if err := uc.userRepo.Update(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to update last login: %w", err)
	}

	// Gerar token JWT
	token, err := uc.authSvc.GenerateToken(user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return &AuthenticateOutput{
		Token:    token,
		UserID:   user.ID,
		Email:    user.Email,
		Role:     string(user.Role),
		Username: user.Username,
	}, nil
}

// AuthenticateOutput representa a saída da autenticação
type AuthenticateOutput struct {
	Token    string
	UserID   string
	Email    string
	Role     string
	Username string
}
