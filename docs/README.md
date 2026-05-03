# 🚀 Sistema de Auditoria Distribuída para Transações Financeiras

## 📖 Índice
- [Visão Geral](#visão-geral)
- [Começando](#começando)
- [Estrutura do Projeto](#estrutura-do-projeto)
- [Arquitetura](#arquitetura)
- [API Endpoints](#api-endpoints)
- [Desenvolvimento](#desenvolvimento)
- [Deploy](#deploy)

---

## 📌 Visão Geral

**Sistema de Auditoria Distribuída** processa transações financeiras de forma assíncrona com registro em blockchain privada, garantindo compliance financeiro, rastreabilidade completa e auditoria confiável.

### ✨ Características

✅ Processamento assíncrono com workers  
✅ Registro imutável em blockchain  
✅ Autenticação JWT  
✅ Arquitetura de microserviços  
✅ Clean Architecture  
✅ Fila distribuída (RabbitMQ/Redis)  
✅ PostgreSQL + SQLite  
✅ Docker Compose pronto  
✅ Logging centralizado  
✅ Injeção de dependências  

---

## 🚀 Começando

### Pré-requisitos

- Go 1.22+
- Docker & Docker Compose
- PostgreSQL (ou SQLite para dev)
- RabbitMQ (opcional, pode usar in-memory)

### Instalação Rápida

```bash
# 1. Clonar repositório
git clone https://github.com/Jeielsantosdev/Work_processamento_assincrono
cd Work_processamento_assincrono

# 2. Copiar .env
cp .env.example .env

# 3. Instalar dependências Go
go mod download

# 4. Iniciar com Docker Compose
docker-compose up -d

# 5. API disponível em: http://localhost:8080
```

### Desenvolvimento Local

```bash
# Terminal 1 - API
go run ./cmd/api/main.go

# Terminal 2 - Worker
go run ./cmd/worker/worker.go
```

---

## 📂 Estrutura do Projeto

```
.
├── cmd/
│   ├── api/                 # Aplicação API
│   │   └── main.go
│   └── worker/              # Processador assíncrono
│       └── worker.go
│
├── internal/
│   ├── domain/
│   │   ├── entities/        # Modelos de domínio
│   │   └── interfaces/      # Contratos
│   │
│   ├── usecase/             # Lógica de negócio
│   ├── infra/               # Implementações infraestrutura
│   │   ├── repository/      # Persistência
│   │   ├── http/            # REST API
│   │   ├── database/        # BD
│   │   ├── queue/           # Filas
│   │   ├── blockchain/      # Blockchain
│   │   └── auth/            # JWT
│   │
│   ├── container/           # Injeção dependências
│   └── worker/              # Workers assíncronos
│
├── pkg/
│   ├── logger/              # Logging
│   ├── errors/              # Erros
│   └── crypto/              # Criptografia
│
├── config/                  # Configurações
├── docs/                    # Documentação
├── docker-compose.yml       # Orquestração containers
├── dockerfile               # Build
├── go.mod & go.sum         # Dependências
└── .env                     # Variáveis ambiente
```

---

## 🏛️ Arquitetura

### Clean Architecture em 4 Camadas

```
┌─────────────────────────────┐
│     API REST / Handlers     │ ← Adapters
├─────────────────────────────┤
│      Use Cases / Business   │ ← Lógica de Aplicação
├─────────────────────────────┤
│    Domain Entities & Rules  │ ← Regras de Negócio
├─────────────────────────────┤
│  Database/Queue/Blockchain  │ ← Infraestrutura
└─────────────────────────────┘
```

### Fluxo de Requisição

```
Cliente
  ↓
Handler HTTP
  ↓
Use Case (ValidationValidate + Business Logic)
  ↓
Repository (Save)
  ↓
Database
  ↓
Queue (Publish for async processing)
  ↓
Worker
  ↓
Blockchain Service (Record)
  ↓
Notification Service (Send)
```

Veja [ARCHITECTURE.md](ARCHITECTURE.md) para detalhes completos.

---

## 🔌 API Endpoints

### Autenticação

```bash
# Login
POST /auth/login
{
  "email": "user@example.com",
  "password": "password123"
}

# Response
{
  "token": "eyJhbGciOiJIUzI1NiIs...",
  "user_id": "uuid",
  "email": "user@example.com",
  "role": "CLIENT"
}
```

### Transações

```bash
# Criar transação
POST /transactions
{
  "client_id": "uuid",
  "source_account": "ACC-001",
  "destination_account": "ACC-002",
  "amount": 1000.50,
  "currency": "BRL",
  "description": "Payment for services"
}

# Listar transações do cliente
GET /transactions?client_id=uuid

# Obter transação específica
GET /transactions/:id
```

### Auditoria

```bash
# Listar logs de auditoria
GET /audit/logs

# Auditoria de transação específica
GET /audit/transactions/:id
```

### Health Check

```bash
GET /health
```

---

## 🛠️ Desenvolvimento

### Adicionar Nova Funcionalidade

1. **Criar Entity** em `internal/domain/entities/`
2. **Definir Interface** em `internal/domain/interfaces/`
3. **Criar Use Case** em `internal/usecase/`
4. **Implementar Repository** em `internal/infra/repository/`
5. **Criar Handler** em `internal/infra/http/handlers/`
6. **Registrar em Container** em `internal/container/container.go`
7. **Adicionar Rotas** em `internal/infra/http/routes/`

### Exemplo: Criar Nova Funcionalidade

```go
// 1. Entity
type Order struct {
    ID string
    // fields...
}

// 2. Interface
type OrderRepository interface {
    Save(ctx context.Context, order *Order) error
    FindByID(ctx context.Context, id string) (*Order, error)
}

// 3. Use Case
type CreateOrderUseCase struct {
    repo OrderRepository
}

func (uc *CreateOrderUseCase) Execute(ctx context.Context, input *Input) (*Output, error) {
    // business logic
}

// 4. Repository
type OrderRepositorySQL struct {
    db *sql.DB
}

// 5. Handler
type OrderHandler struct {
    container *Container
}

// 6. Container
orderRepo := repository.NewOrderRepositorySQL(db)
createOrderUC := usecase.NewCreateOrderUseCase(orderRepo)

// 7. Routes
router.POST("/orders", orderHandler.CreateOrder)
```

### Rodando Testes

```bash
# Todos os testes
go test ./...

# Com coverage
go test -cover ./...

# Específico
go test ./internal/usecase/...
```

---

## 🐳 Docker & Deploy

### Docker Compose

```bash
# Iniciar todos os serviços
docker-compose up -d

# Parar serviços
docker-compose down

# Ver logs
docker-compose logs -f api
docker-compose logs -f worker
```

### Serviços

- **API**: `http://localhost:8080`
- **RabbitMQ Management**: `http://localhost:15672`
- **PostgreSQL**: `localhost:5432`
- **Redis**: `localhost:6379`

### Variáveis de Ambiente

Edite `.env` para configurar:

```env
APP_ENV=production
APP_PORT=:8080

# Database
DB_DRIVER=postgres
DB_CONNECTION_STRING=postgres://user:pass@host:5432/db

# Queue
QUEUE_TYPE=rabbitmq
QUEUE_CONNECTION_STRING=amqp://guest:guest@rabbitmq:5672/

# JWT
JWT_SECRET_KEY=your-secret-key-here
JWT_EXPIRES_IN=24

# Logger
LOG_LEVEL=info
```

### Build e Push para Registro

```bash
# Build
docker build -t auditoria:1.0 .

# Tag
docker tag auditoria:1.0 seu-registry/auditoria:1.0

# Push
docker push seu-registry/auditoria:1.0
```

---

## 📊 Monitoramento & Logging

### Logs Estruturados

```
[2024-04-28] INFO: Transaction created | tx_id=uuid | amount=1000
[2024-04-28] ERROR: Failed to process | tx_id=uuid | error=insufficient_balance
```

### Métricas

Adicione Prometheus para métricas:

```go
import "github.com/prometheus/client_golang/prometheus"

transactionCounter := prometheus.NewCounter(
    prometheus.CounterOpts{
        Name: "transactions_total",
    },
)
```

---

## 🔐 Segurança

- ✅ JWT para autenticação
- ✅ Bcrypt para hash de senhas
- ✅ Validação em camadas
- ✅ TLS para comunicação
- ⚠️ TODO: Rate limiting
- ⚠️ TODO: CORS
- ⚠️ TODO: HTTPS
- ⚠️ TODO: Secrets management

---

## 📚 Tecnologias

| Layer | Tecnologia |
|-------|-----------|
| **Language** | Go 1.22 |
| **API** | HTTP REST |
| **Database** | PostgreSQL / SQLite |
| **Queue** | RabbitMQ / Redis |
| **Auth** | JWT |
| **Crypto** | Bcrypt, SHA256 |
| **Logging** | Logrus |
| **Container** | Docker |
| **Orchestration** | Docker Compose |

---

## 🤝 Contribuindo

1. Fork o projeto
2. Crie uma branch (`git checkout -b feature/AmazingFeature`)
3. Commit suas mudanças (`git commit -m 'Add AmazingFeature'`)
4. Push para a branch (`git push origin feature/AmazingFeature`)
5. Abra um Pull Request

---

## 📄 Licença

Distribuído sob a licença MIT. Veja `LICENSE` para mais informações.

---

## 👤 Autor

**Jeiel Santos**
- GitHub: [@Jeielsantosdev](https://github.com/Jeielsantosdev)

---

## 📞 Suporte

Para dúvidas ou sugestões:
- 📧 Email: jeiel@example.com
- 💬 Issues: [GitHub Issues](https://github.com/Jeielsantosdev/Work_processamento_assincrono/issues)

---

## 🙏 Agradecimentos

- Clean Architecture - Robert C. Martin
- Microservices Architecture - Sam Newman
- Go Community


---

## 9. Diagramas UML

### 9.1. Diagrama de Casos de Uso

```mermaid
usecaseDiagram
  actor Cliente
  actor Auditor
  actor Administrador

  usecase UC01 as "Enviar Transação"
  usecase UC02 as "Receber Status"
  usecase UC03 as "Auditar Histórico"
  usecase UC04 as "Gerenciar Sistema"

  Cliente --> UC01
  Cliente --> UC02
  Auditor --> UC03
  Administrador --> UC04
```

## 9.2. Diagrama de Classes

    ````classDiagram

class Usuario {
+id: int
+nome: string
+tipo: string
+credenciais: string
+autenticar()
}

class Transacao {
+id: int
+origem: string
+destino: string
+valor: float
+status: string
+timestamp: datetime
+validar()
+processar()
}

class Bloco {
+hash: string
+transacoes: list<Transacao>
+hashAnterior: string
+timestamp: datetime
+assinar()
}

class Notificacao {
+id: int
+transacao_id: int
+mensagem: string
+status: string
+timestamp: datetime
+enviar()
}

Usuario "1" --> "N" Transacao
Transacao "1" --> "1" Bloco
Transacao "1" --> "N" Notificacao

``````



## 9.3. Diagrama de Sequência – Fluxo Assíncrono

    ````sequenceDiagram
    Cliente ->> TransactionService: Envia Transação (gRPC)
    TransactionService ->> TransactionService: Valida Transação
    TransactionService ->> TransactionService: Gera ID Único
    TransactionService ->> Queue: Envia Transação
    TransactionService -->> Cliente: Retorna ID

        Queue ->> WorkerPool: Processa Transação
        WorkerPool ->> WorkerPool: Valida Saldo/Regras
        WorkerPool ->> BlockchainWriter: Prepara dados
        BlockchainWriter ->> Blockchain: Grava Bloco
        BlockchainWriter -->> WorkerPool: Confirmação
        WorkerPool ->> NotificationService: Envia Status
        NotificationService -->> Cliente: Notificação
    `````
``````
