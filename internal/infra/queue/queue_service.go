package queue

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Jeielsantosdev/Work_processamento_assincrono/internal/domain/entities"
)

// InMemoryQueueService implementa um serviço de fila em memória (para desenvolvimento)
type InMemoryQueueService struct {
	transactions chan *entities.Transaction
	maxSize      int
}

// NewInMemoryQueueService cria uma nova instância da fila em memória
func NewInMemoryQueueService(maxSize int) *InMemoryQueueService {
	if maxSize <= 0 {
		maxSize = 1000
	}
	return &InMemoryQueueService{
		transactions: make(chan *entities.Transaction, maxSize),
		maxSize:      maxSize,
	}
}

// PublishTransaction publica uma transação na fila
func (q *InMemoryQueueService) PublishTransaction(ctx context.Context, tx *entities.Transaction) error {
	select {
	case q.transactions <- tx:
		fmt.Printf("[QUEUE] Published transaction: %s\n", tx.ID)
		return nil
	case <-ctx.Done():
		return ctx.Err()
	default:
		return fmt.Errorf("queue is full")
	}
}

// ConsumeTransactions consome transações da fila
func (q *InMemoryQueueService) ConsumeTransactions(ctx context.Context) (<-chan *entities.Transaction, error) {
	return q.transactions, nil
}

// AcknowledgeMessage reconhece um mensagem (não operacional na fila em memória)
func (q *InMemoryQueueService) AcknowledgeMessage(ctx context.Context, messageID string) error {
	return nil
}

// RabbitMQQueueService implementa o serviço de fila com RabbitMQ (estrutura base)
type RabbitMQQueueService struct {
	connectionString string
	queueName        string
	// connection *amqp.Connection
	// channel    *amqp.Channel
}

// NewRabbitMQQueueService cria uma nova instância do serviço RabbitMQ
func NewRabbitMQQueueService(connectionString, queueName string) *RabbitMQQueueService {
	return &RabbitMQQueueService{
		connectionString: connectionString,
		queueName:        queueName,
	}
}

// PublishTransaction publica uma transação na fila RabbitMQ
func (q *RabbitMQQueueService) PublishTransaction(ctx context.Context, tx *entities.Transaction) error {
	// TODO: Implementar integração com RabbitMQ
	data, err := json.Marshal(tx)
	if err != nil {
		return err
	}

	fmt.Printf("[RABBITMQ] Published transaction: %s with data: %s\n", tx.ID, string(data))
	return nil
}

// ConsumeTransactions consome transações da fila RabbitMQ
func (q *RabbitMQQueueService) ConsumeTransactions(ctx context.Context) (<-chan *entities.Transaction, error) {
	// TODO: Implementar integração com RabbitMQ
	ch := make(chan *entities.Transaction)
	return ch, nil
}

// AcknowledgeMessage reconhece uma mensagem no RabbitMQ
func (q *RabbitMQQueueService) AcknowledgeMessage(ctx context.Context, messageID string) error {
	// TODO: Implementar integração com RabbitMQ
	return nil
}
