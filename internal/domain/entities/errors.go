package entities

import "errors"

// Erros de validação de domínio
var (
	ErrInvalidClientID     = errors.New("client ID is invalid or empty")
	ErrInvalidAccount      = errors.New("source or destination account is invalid or empty")
	ErrInvalidAmount       = errors.New("amount must be greater than zero")
	ErrSameAccount         = errors.New("source and destination accounts cannot be the same")
	ErrInsufficientBalance = errors.New("insufficient balance for transaction")
	ErrTransactionNotFound = errors.New("transaction not found")
	ErrUserNotFound        = errors.New("user not found")
	ErrUnauthorized        = errors.New("user is not authorized to perform this action")
	ErrInvalidCredentials  = errors.New("invalid credentials provided")
	ErrUserAlreadyExists   = errors.New("user already exists")
	ErrBlockchainError     = errors.New("error recording transaction in blockchain")
	ErrDatabaseError       = errors.New("database operation failed")
	ErrQueueError          = errors.New("queue operation failed")
	ErrInternalError       = errors.New("internal server error")
	ErrNotificationFailed  = errors.New("failed to send notification")
)
