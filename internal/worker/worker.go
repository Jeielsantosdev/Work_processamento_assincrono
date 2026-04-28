package worker

import (
	"context"
	"fmt"
	"sync"

	"github.com/Jeielsantosdev/Work_processamento_assincrono/internal/container"
	"github.com/Jeielsantosdev/Work_processamento_assincrono/internal/domain/entities"
	"github.com/Jeielsantosdev/Work_processamento_assincrono/pkg/logger"
)

// Worker processa transações de forma assíncrona
type Worker struct {
	container   *container.Container
	logger      logger.Logger
	workerCount int
	wg          sync.WaitGroup
	ctx         context.Context
	cancel      context.CancelFunc
}

// NewWorker cria um novo worker
func NewWorker(c *container.Container, workerCount int) *Worker {
	ctx, cancel := context.WithCancel(context.Background())
	return &Worker{
		container:   c,
		logger:      c.Logger,
		workerCount: workerCount,
		ctx:         ctx,
		cancel:      cancel,
	}
}

// Start inicia os workers
func (w *Worker) Start() error {
	w.logger.Info("Starting workers", "count", w.workerCount)

	// Obter canal de transações da fila
	transactionChan, err := w.container.QueueService.ConsumeTransactions(w.ctx)
	if err != nil {
		return fmt.Errorf("failed to get transaction channel: %w", err)
	}

	// Iniciar workers
	for i := 0; i < w.workerCount; i++ {
		w.wg.Add(1)
		go w.processTransactionsWorker(i+1, transactionChan)
	}

	w.logger.Info("Workers started successfully")
	return nil
}

// processTransactionsWorker processa transações de um worker
func (w *Worker) processTransactionsWorker(workerID int, transactionChan <-chan *entities.Transaction) {
	defer w.wg.Done()

	for tx := range transactionChan {
		select {
		case <-w.ctx.Done():
			w.logger.Info(fmt.Sprintf("Worker %d stopping", workerID))
			return
		default:
			w.logger.Info(fmt.Sprintf("Worker %d processing transaction %s", workerID, tx.ID))

			// TODO: Chamar ProcessTransactionUC.Execute()
			// Por enquanto, apenas log
			tx.MarkAsCompleted(tx.ID)

			w.logger.Info(fmt.Sprintf("Worker %d completed transaction %s", workerID, tx.ID))
		}
	}
}

// Stop para os workers
func (w *Worker) Stop() {
	w.logger.Info("Stopping workers")
	w.cancel()
	w.wg.Wait()
	w.logger.Info("Workers stopped")
}

// WorkerPool gerencia múltiplos workers
type WorkerPool struct {
	workers []*Worker
	logger  logger.Logger
}

// NewWorkerPool cria um novo pool de workers
func NewWorkerPool(container *container.Container, poolSize int) *WorkerPool {
	workers := make([]*Worker, poolSize)
	for i := 0; i < poolSize; i++ {
		workers[i] = NewWorker(container, 1)
	}

	return &WorkerPool{
		workers: workers,
		logger:  container.Logger,
	}
}

// Start inicia todos os workers do pool
func (wp *WorkerPool) Start() error {
	for i, worker := range wp.workers {
		if err := worker.Start(); err != nil {
			wp.logger.Error("Failed to start worker", err, "worker_id", i)
			return err
		}
	}
	return nil
}

// Stop para todos os workers do pool
func (wp *WorkerPool) Stop() {
	for _, worker := range wp.workers {
		worker.Stop()
	}
}
