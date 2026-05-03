package worker

import (
	"context"
	"testing"
	"time"

	"github.com/Jeielsantosdev/Work_processamento_assincrono/internal/container"
	"github.com/Jeielsantosdev/Work_processamento_assincrono/internal/domain/entities"
	qinfra "github.com/Jeielsantosdev/Work_processamento_assincrono/internal/infra/queue"
	"github.com/Jeielsantosdev/Work_processamento_assincrono/pkg/logger"
)

func TestWorkerProcessesInMemoryQueue(t *testing.T) {
	cfgLogger := logger.NewLogrusLogger("debug")
	c := &container.Container{
		Logger:       cfgLogger,
		QueueService: qinfra.NewInMemoryQueueService(10),
		// ProcessTransactionUC left nil to exercise fallback path
	}

	w := NewWorker(c, 1)
	if err := w.Start(); err != nil {
		t.Fatalf("failed to start worker: %v", err)
	}

	tx := &entities.Transaction{ID: "tx-test", ClientID: "cli-1", Amount: 1.0, Status: entities.StatusPending}
	// publish
	if err := c.QueueService.PublishTransaction(context.Background(), tx); err != nil {
		t.Fatalf("failed to publish tx: %v", err)
	}

	// aguardar processamento
	time.Sleep(200 * time.Millisecond)

	if tx.Status != entities.StatusCompleted {
		t.Fatalf("expected transaction to be completed by worker, got %s", tx.Status)
	}

	w.Stop()
}
