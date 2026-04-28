# 📋 Documentação da Arquitetura

## 🏛️ Arquitetura Geral

O sistema é construído seguindo os princípios de **Clean Architecture** e **Microserviços**, dividido em:

### 1. **API Service** (`cmd/api/main.go`)
- Recebe requisições HTTP/REST
- Valida e encaminha para casos de uso
- Retorna respostas estruturadas
- Porta: `8080`

### 2. **Worker Service** (`cmd/worker/worker.go`)
- Processa transações de forma assíncrona
- Consome mensagens da fila
- Coordena com blockchain e notificação
- Pool de workers configurável

## 🎯 Camadas (Clean Architecture)

### **Entities** (`internal/domain/entities/`)
- Regras de negócio críticas
- Sem dependências externas
- Validações intrínsecas
- Exemplo: `Transaction`, `User`, `BlockRecord`

### **Use Cases** (`internal/usecase/`)
- Lógica de aplicação
- Independentes de frameworks
- Orquestram domínio e interfaces
- Exemplos:
  - `CreateTransactionUseCase`
  - `ProcessTransactionUseCase`
  - `AuthenticateUserUseCase`

### **Interfaces** (`internal/domain/interfaces/`)
- Contratos entre camadas
- Define a comunicação
- Exemplos:
  - `TransactionRepository`
  - `BlockchainService`
  - `QueueService`

### **Adapters/Infra** (`internal/infra/`)
- Implementações concretas
- Integração com banco, fila, blockchain
- Handlers HTTP/gRPC
- Subpastas:
  - `repository/` - Persistência
  - `queue/` - Filas (RabbitMQ, Redis)
  - `blockchain/` - Blockchain
  - `auth/` - JWT
  - `http/` - REST API
  - `database/` - Conexão BD

### **Pacotes Utilitários** (`pkg/`)
- `logger/` - Logging centralizado
- `errors/` - Tratamento de erros
- `crypto/` - Criptografia

## 🔌 Estrutura de Diretórios

```
internal/
├── domain/
│   ├── entities/          # Modelos de negócio
│   ├── interfaces/        # Contratos
│
├── usecase/               # Casos de uso
│
├── infra/
│   ├── repository/        # Implementações de persistência
│   ├── http/
│   │   ├── handlers/      # Handlers REST
│   │   └── routes/        # Roteamento
│   ├── grpc/              # Handlers gRPC
│   ├── database/          # Conexão e migrations
│   ├── queue/             # Implementações de fila
│   ├── blockchain/        # Implementações blockchain
│   ├── auth/              # JWT
│   ├── crypto/            # Criptografia
│   └── notification/      # Notificações
│
├── container/             # Injeção de dependências
├── worker/                # Processador assíncrono
└── models/                # DTOs (Data Transfer Objects)

cmd/
├── api/                   # Ponto de entrada API
└── worker/                # Ponto de entrada Worker
```

## 🔄 Fluxo de Requisição

```
Cliente HTTP
    ↓
Handler REST (routes/handlers)
    ↓
Use Case (lógica de negócio)
    ↓
Repository (persistência)
    ↓
Database
```

## ⚙️ Fluxo de Processamento Assíncrono

```
API: Criar Transação
    ↓
[Fila de Transações]
    ↓
Worker: Processar
    ↓
Use Case: ProcessTransactionUseCase
    ↓
Blockchain Service: Registrar
    ↓
Notification Service: Notificar Cliente
    ↓
Database: Atualizar Status
```

## 📊 Dependências Injetadas (Container)

O container gerencia todas as dependências:

```go
Container {
    // Repositórios
    TransactionRepo
    UserRepo
    NotificationRepo
    
    // Serviços
    BlockchainService
    NotificationService
    QueueService
    
    // Use Cases
    CreateTransactionUC
    ProcessTransactionUC
}
```

## 🛡️ Princípios SOLID Aplicados

### **S**ingle Responsibility
- Cada classe tem uma responsabilidade
- `TransactionRepository` - apenas persistência

### **O**pen/Closed
- Aberto para extensão via interfaces
- `BlockchainService` pode ser substituído

### **L**iskov Substitution
- Implementações de interfaces são intercambiáveis
- `SimpleBlockchain` ↔ `EthereumBlockchain`

### **I**nterface Segregation
- Interfaces específicas por domínio
- `TransactionRepository` não sabe de `UserRepository`

### **D**ependency Inversion
- Depende de abstrações, não de implementações
- Use cases usam interfaces, não implementações diretas

## 🔐 Segurança

- **JWT** para autenticação stateless
- **Bcrypt** para hash de senha
- **Validação** em múltiplas camadas
- **TLS** recomendado em produção
- **Rate limiting** (implementar)
- **CORS** (implementar)

## 📈 Escalabilidade

- **Stateless API** - escala horizontalmente
- **Worker Pool** - processa paralelo
- **Database Connection Pool** - otimizado
- **Fila distribuída** - RabbitMQ/Redis
- **Múltiplas instâncias** - via Docker Compose

## 🔧 Configuração

Todas as configurações em `config/config.go` e `.env`:

```
APP_ENV=development
DB_DRIVER=postgres
QUEUE_TYPE=rabbitmq
LOG_LEVEL=info
```

## 📝 Extensões Futuras

1. **gRPC** - Substituir HTTP por gRPC
2. **Rate Limiting** - Implementar via middleware
3. **Caching** - Redis para dados frequentes
4. **Observabilidade** - Prometheus + Grafana
5. **Testes** - Unit, Integration, E2E
6. **CI/CD** - GitHub Actions
7. **Documentação API** - Swagger/OpenAPI
