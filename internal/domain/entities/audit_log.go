package entities

import (
	"time"

	"github.com/google/uuid"
)

// AuditLog representa um registro de auditoria do sistema
type AuditLog struct {
	ID            string
	TransactionID string
	Action        string
	ActorID       string
	ActorRole     UserRole
	Details       map[string]interface{}
	IPAddress     string
	Timestamp     time.Time
	Status        string
	ErrorMessage  string
}

// NewAuditLog cria um novo registro de auditoria
func NewAuditLog(txID, action, actorID string, role UserRole, details map[string]interface{}, ipAddress string) *AuditLog {
	return &AuditLog{
		ID:            uuid.New().String(),
		TransactionID: txID,
		Action:        action,
		ActorID:       actorID,
		ActorRole:     role,
		Details:       details,
		IPAddress:     ipAddress,
		Timestamp:     time.Now(),
		Status:        "SUCCESS",
	}
}

// MarkAsError marca o log como erro
func (a *AuditLog) MarkAsError(errMsg string) {
	a.Status = "ERROR"
	a.ErrorMessage = errMsg
}
