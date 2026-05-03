package entities

import (
	"time"

	"github.com/google/uuid"
)

// NotificationType representa o tipo de notificação
type NotificationType string

const (
	NotificationTypeTransactionStatus NotificationType = "TRANSACTION_STATUS"
	NotificationTypeAlert             NotificationType = "ALERT"
	NotificationTypeReport            NotificationType = "REPORT"
)

// NotificationStatus representa o status de entrega da notificação
type NotificationStatus string

const (
	NotificationStatusPending NotificationStatus = "PENDING"
	NotificationStatusSent    NotificationStatus = "SENT"
	NotificationStatusFailed  NotificationStatus = "FAILED"
)

// Notification representa uma notificação para cliente
type Notification struct {
	ID            string
	ClientID      string
	Type          NotificationType
	Title         string
	Message       string
	Data          map[string]interface{}
	Status        NotificationStatus
	Channel       string // email, sms, push
	SentAt        *time.Time
	FailureReason string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

// NewNotification cria uma nova notificação
func NewNotification(clientID string, notifType NotificationType, title, message string, data map[string]interface{}, channel string) *Notification {
	return &Notification{
		ID:        uuid.New().String(),
		ClientID:  clientID,
		Type:      notifType,
		Title:     title,
		Message:   message,
		Data:      data,
		Status:    NotificationStatusPending,
		Channel:   channel,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// MarkAsSent marca a notificação como enviada
func (n *Notification) MarkAsSent() {
	n.Status = NotificationStatusSent
	now := time.Now()
	n.SentAt = &now
	n.UpdatedAt = now
}

// MarkAsFailed marca a notificação como falhada
func (n *Notification) MarkAsFailed(reason string) {
	n.Status = NotificationStatusFailed
	n.FailureReason = reason
	n.UpdatedAt = time.Now()
}
