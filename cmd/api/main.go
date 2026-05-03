package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/Jeielsantosdev/Work_processamento_assincrono/config"
	"github.com/Jeielsantosdev/Work_processamento_assincrono/internal/container"
	"github.com/Jeielsantosdev/Work_processamento_assincrono/internal/infra/http/routes"
	"github.com/Jeielsantosdev/Work_processamento_assincrono/internal/worker"
)

func main() {
	// Carregar configuração
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load config: %v\n", err)
		os.Exit(1)
	}

	// Validar configuração
	if err := cfg.Validate(); err != nil {
		fmt.Fprintf(os.Stderr, "Invalid configuration: %v\n", err)
		os.Exit(1)
	}

	// Criar container de dependências
	c, err := container.New(*cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create container: %v\n", err)
		os.Exit(1)
	}
	defer c.Close()

	c.Logger.Info("Application started", "name", cfg.App.Name, "version", cfg.App.Version)

	// Configurar rotas
	router := routes.SetupRoutes(c)

	// Iniciar worker pool
	workerPool := worker.NewWorkerPool(c, 4) // 4 workers
	if err := workerPool.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to start worker pool: %v\n", err)
		os.Exit(1)
	}

	// Servidor HTTP
	server := &http.Server{
		Addr:    cfg.App.Port,
		Handler: router,
	}

	// Iniciar servidor em goroutine
	go func() {
		c.Logger.Info("Starting HTTP server", "addr", cfg.App.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			c.Logger.Error("Server error", err)
		}
	}()

	// Aguardar sinal de interrupção
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	<-sigChan

	c.Logger.Info("Received shutdown signal")

	// Parar workers
	workerPool.Stop()

	// Shutdown gracioso do servidor
	if err := server.Close(); err != nil {
		c.Logger.Error("Error during server shutdown", err)
	}

	c.Logger.Info("Application stopped")
}
