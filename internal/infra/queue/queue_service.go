package queue

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"

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

// RabbitMQQueueService implementa o serviço de fila com RabbitMQ
type RabbitMQQueueService struct {
	connectionString string
	queueName        string
	dlqName          string

	conn    *amqp.Connection
	channel *amqp.Channel

	deliveries map[string]amqp.Delivery
	mu         sync.Mutex
}

// NewRabbitMQQueueService cria uma nova instância do serviço RabbitMQ
func NewRabbitMQQueueService(connectionString, queueName string) *RabbitMQQueueService {
	return &RabbitMQQueueService{
		connectionString: connectionString,
		queueName:        queueName,
		dlqName:          queueName + ".dlq",
		deliveries:       make(map[string]amqp.Delivery),
	}
}

// ensureConnection estabelece conexão e canal com RabbitMQ
func (q *RabbitMQQueueService) ensureConnection() error {
	if q.conn != nil && q.channel != nil {
		return nil
	}

	conn, err := amqp.Dial(q.connectionString)
	if err != nil {
		return fmt.Errorf("failed to dial rabbitmq: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return fmt.Errorf("failed to open channel: %w", err)
	}

	// Declarar fila principal e DLQ
	_, err = ch.QueueDeclare(
		q.queueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		ch.Close()
		conn.Close()
		return fmt.Errorf("failed to declare queue: %w", err)
	}

	_, err = ch.QueueDeclare(
		q.dlqName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		ch.Close()
		conn.Close()
		return fmt.Errorf("failed to declare dlq: %w", err)
	}

	q.conn = conn
	q.channel = ch
	return nil
}

// PublishTransaction publica uma transação na fila RabbitMQ
func (q *RabbitMQQueueService) PublishTransaction(ctx context.Context, tx *entities.Transaction) error {
	if err := q.ensureConnection(); err != nil {
		return err
	}

	data, err := json.Marshal(tx)
	if err != nil {
		return err
	}

	headers := amqp.Table{}
	headers["x-retry-count"] = tx.RetryCount

	err = q.channel.PublishWithContext(ctx,
		"", // default exchange
		q.queueName,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        data,
			MessageId:   tx.ID,
			Headers:     headers,
			Timestamp:   time.Now(),
		},
	)
	if err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}

	fmt.Printf("[RABBITMQ] Published transaction: %s\n", tx.ID)
	return nil
}

// ConsumeTransactions consome transações da fila RabbitMQ
func (q *RabbitMQQueueService) ConsumeTransactions(ctx context.Context) (<-chan *entities.Transaction, error) {
	if err := q.ensureConnection(); err != nil {
		return nil, err
	}

	msgs, err := q.channel.Consume(
		q.queueName,
		"",
		false, // autoAck false: we'll ack manually
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to start consume: %w", err)
	}

	out := make(chan *entities.Transaction)

	go func() {
		defer close(out)
		for d := range msgs {
			var tx entities.Transaction
			if err := json.Unmarshal(d.Body, &tx); err != nil {
				// Malformed message: move to DLQ and ack
				_ = q.channel.Publish(
					"",
					q.dlqName,
					false,
					false,
					amqp.Publishing{ContentType: "text/plain", Body: []byte("invalid payload")},
				)
				d.Ack(false)
				continue
			}

			// Set retry count from headers if present
			if v, ok := d.Headers["x-retry-count"]; ok {
				switch val := v.(type) {
				case int32:
					tx.RetryCount = int(val)
				case int64:
					tx.RetryCount = int(val)
				case int:
					tx.RetryCount = val
				case float64:
					tx.RetryCount = int(val)
				}
			}

			// If exceeded retries, move to DLQ
			if tx.RetryCount >= tx.MaxRetries {
				_ = q.channel.Publish(
					"",
					q.dlqName,
					false,
					false,
					amqp.Publishing{ContentType: "application/json", Body: d.Body},
				)
				d.Ack(false)
				continue
			}

			// Store delivery to allow ack/nack later by message id
			q.mu.Lock()
			q.deliveries[tx.ID] = d
			q.mu.Unlock()

			select {
			case out <- &tx:
				// forwarded to worker
			case <-ctx.Done():
				// Context canceled: nack and exit
				d.Nack(false, true)
				q.mu.Lock()
				delete(q.deliveries, tx.ID)
				q.mu.Unlock()
				return
			}
		}
	}()

	return out, nil
}

// AcknowledgeMessage reconhece uma mensagem no RabbitMQ por messageID (tx.ID)
func (q *RabbitMQQueueService) AcknowledgeMessage(ctx context.Context, messageID string) error {
	q.mu.Lock()
	d, ok := q.deliveries[messageID]
	if !ok {
		q.mu.Unlock()
		return errors.New("message not found")
	}
	delete(q.deliveries, messageID)
	q.mu.Unlock()

	if err := d.Ack(false); err != nil {
		return fmt.Errorf("failed to ack message: %w", err)
	}
	return nil
}
