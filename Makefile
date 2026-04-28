.PHONY: help build run test dev docker-up docker-down clean lint format install-tools

help:
	@echo "╔══════════════════════════════════════════════════════════════╗"
	@echo "║  Sistema de Auditoria Distribuída - Makefile                ║"
	@echo "╚══════════════════════════════════════════════════════════════╝"
	@echo ""
	@echo "Desenvolvimento:"
	@echo "  make dev              - Inicia API e Worker em desenvolvimento"
	@echo "  make install-deps     - Instala dependências Go"
	@echo "  make test             - Roda testes"
	@echo "  make lint             - Verifica linting (requer golangci-lint)"
	@echo "  make format           - Formata código (gofmt)"
	@echo ""
	@echo "Build & Deploy:"
	@echo "  make build            - Compila binários (api e worker)"
	@echo "  make build-api        - Compila apenas API"
	@echo "  make build-worker     - Compila apenas Worker"
	@echo ""
	@echo "Docker:"
	@echo "  make docker-build     - Build imagens Docker"
	@echo "  make docker-up        - Inicia stack completa"
	@echo "  make docker-down      - Para stack"
	@echo "  make docker-logs      - Mostra logs"
	@echo "  make docker-clean     - Remove containers e volumes"
	@echo ""
	@echo "Utilitários:"
	@echo "  make clean            - Remove arquivos gerados"
	@echo "  make mod-tidy         - Limpa dependências"
	@echo "  make help             - Mostra esta ajuda"

# Desenvolvimento
dev:
	@echo "🚀 Iniciando em modo desenvolvimento..."
	@echo "Terminal 1: API"
	@echo "Terminal 2: Worker"
	@tmux new-session -d -s dev "make run-api; bash" && tmux new-window -t dev -n worker "make run-worker; bash"

run-api:
	@echo "▶️  Iniciando API..."
	go run ./cmd/api/main.go

run-worker:
	@echo "▶️  Iniciando Worker..."
	go run ./cmd/worker/worker.go

install-deps:
	@echo "📦 Instalando dependências..."
	go mod download
	go mod tidy

test:
	@echo "🧪 Rodando testes..."
	go test -v -cover ./...

test-coverage:
	@echo "📊 Gerando coverage..."
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "✅ Coverage report: coverage.html"

lint:
	@echo "🔍 Verificando linting..."
	golangci-lint run ./...

format:
	@echo "🎨 Formatando código..."
	gofmt -s -w .
	go mod tidy

# Build
build: build-api build-worker
	@echo "✅ Build completo!"

build-api:
	@echo "🔨 Building API..."
	CGO_ENABLED=1 go build -o bin/api ./cmd/api

build-worker:
	@echo "🔨 Building Worker..."
	CGO_ENABLED=1 go build -o bin/worker ./cmd/worker

# Docker
docker-build:
	@echo "🐳 Building Docker images..."
	docker build -t auditoria:api --target api .
	docker build -t auditoria:worker --target worker .

docker-up:
	@echo "🚀 Iniciando stack Docker..."
	docker-compose up -d
	@echo "✅ Stack iniciada!"
	@echo "📍 API: http://localhost:8080"
	@echo "📍 RabbitMQ: http://localhost:15672"
	@echo "📍 PostgreSQL: localhost:5432"

docker-down:
	@echo "🛑 Parando stack Docker..."
	docker-compose down

docker-logs:
	@echo "📋 Logs em tempo real..."
	docker-compose logs -f

docker-logs-api:
	docker-compose logs -f api

docker-logs-worker:
	docker-compose logs -f worker

docker-clean:
	@echo "🧹 Limpando Docker..."
	docker-compose down -v
	docker system prune -f

# Utilitários
clean:
	@echo "🧹 Limpando..."
	rm -rf bin/
	rm -f coverage.out coverage.html
	rm -f *.db

mod-tidy:
	@echo "🔄 Limpando dependências..."
	go mod tidy

install-tools:
	@echo "🛠️  Instalando ferramentas de desenvolvimento..."
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@echo "✅ Ferramentas instaladas!"

# Info
info:
	@echo "╔══════════════════════════════════════════════════════════════╗"
	@echo "║              Informações do Projeto                          ║"
	@echo "╚══════════════════════════════════════════════════════════════╝"
	@echo "Go Version: $(shell go version)"
	@echo "Go Path: $(shell go env GOPATH)"
	@echo "Current Dir: $(shell pwd)"
	@echo "Files Changed: $(shell git status --short | wc -l)"
	@echo ""
	@echo "Módulo: $(shell grep module go.mod | awk '{print $$2}')"
	@echo "Go Version Required: $(shell grep '^go ' go.mod)"
	@echo ""
	@echo "Estrutura:"
	@echo "  - Entidades: $(shell find internal/domain/entities -type f -name '*.go' | wc -l) arquivos"
	@echo "  - Use Cases: $(shell find internal/usecase -type f -name '*.go' | wc -l) arquivos"
	@echo "  - Handlers: $(shell find internal/infra/http/handlers -type f -name '*.go' | wc -l) arquivos"

.DEFAULT_GOAL := help
