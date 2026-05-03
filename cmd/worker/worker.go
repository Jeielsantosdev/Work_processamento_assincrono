package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/Jeielsantosdev/Work_processamento_assincrono/config"
	"github.com/Jeielsantosdev/Work_processamento_assincrono/internal/container"
	"github.com/Jeielsantosdev/Work_processamento_assincrono/internal/worker"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load config: %v\n", err)
		os.Exit(1)
	}

	if err := cfg.Validate(); err != nil {
		fmt.Fprintf(os.Stderr, "Invalid configuration: %v\n", err)
		os.Exit(1)
	}

	c, err := container.New(*cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create container: %v\n", err)
		os.Exit(1)
	}
	defer c.Close()

	c.Logger.Info("Worker application started", "name", cfg.App.Name)

	// Criar e iniciar pool de workers
	pool := worker.NewWorkerPool(c, 4)
	if err := pool.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to start worker pool: %v\n", err)
		os.Exit(1)
	}

	// Aguardar sinal de término
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	c.Logger.Info("Shutting down worker pool")
	pool.Stop()
	c.Logger.Info("Worker application stopped")
}
