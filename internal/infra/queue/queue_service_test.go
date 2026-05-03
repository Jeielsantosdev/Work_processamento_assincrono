package queue

import (
	"context"
	"testing"
	"time"

	"github.com/Jeielsantosdev/Work_processamento_assincrono/internal/domain/entities"
)

func TestInMemoryPublishConsume(t *testing.T) {
	q := NewInMemoryQueueService(10)
	ctx := context.Background()

	tx := &entities.Transaction{ID: "tx-1", ClientID: "client-1", Amount: 10.0, Status: entities.StatusPending}

	if err := q.PublishTransaction(ctx, tx); err != nil {
		t.Fatalf("publish failed: %v", err)
	}

	ch, err := q.ConsumeTransactions(ctx)
	if err != nil {
		t.Fatalf("consume failed: %v", err)
	}

	select {
	case got := <-ch:
		if got.ID != tx.ID {
			t.Fatalf("expected tx id %s got %s", tx.ID, got.ID)
		}
	case <-time.After(time.Second):
		t.Fatalf("timed out waiting for message")
	}
}
