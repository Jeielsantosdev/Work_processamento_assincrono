# 🚀 Quick Start Guide

## Inicialização Rápida (5 minutos)

### 1. Clone e Configure

```bash
cd /home/jeiel-dev/Astreus_dev/Work_processamento_assincrono

# Copiar variáveis de ambiente
cp .env.example .env
```

### 2. Com Docker Compose (Recomendado)

```bash
# Iniciar todos os serviços
make docker-up

# Verificar se tudo está rodando
docker-compose ps

# Ver logs em tempo real
make docker-logs
```

**Endpoints:**
- 🔗 API: `http://localhost:8080`
- 🐰 RabbitMQ: `http://localhost:15672` (guest/guest)
- 🐘 PostgreSQL: `localhost:5432`
- 🔴 Redis: `localhost:6379`

### 3. Desenvolvimento Local

```bash
# Terminal 1 - API
make run-api

# Terminal 2 - Worker
make run-worker
```

---

## 📝 Exemplos de API

### Criar Transação

```bash
curl -X POST http://localhost:8080/transactions \
  -H "Content-Type: application/json" \
  -d '{
    "client_id": "client-123",
    "source_account": "ACC-001",
    "destination_account": "ACC-002",
    "amount": 1500.00,
    "currency": "BRL",
    "description": "Pagamento serviço"
  }'
```

**Resposta:**
```json
{
  "transaction_id": "txn-uuid",
  "status": "PENDING",
  "created_at": "2024-04-28T10:30:00Z"
}
```

### Listar Transações

```bash
curl http://localhost:8080/transactions?client_id=client-123
```

### Health Check

```bash
curl http://localhost:8080/health
```

---

## 🏗️ Estrutura Clean Architecture

```
Domain (Entidades, Regras)
  ↑
Application (Use Cases, Interfaces)
  ↑
Adapters (HTTP, DB, Queue)
  ↑
Infrastructure (Docker, Config)
```

---

## 📊 Fluxo de Processamento

1. **Cliente envia transação** → API REST
2. **API valida** → Use Case
3. **Use Case salva** → Repository → Database
4. **Use Case publica** → Queue
5. **Worker consome** → Blockchain Service
6. **Blockchain registra** → Imutável
7. **Notificação enviada** → Email/SMS/Push
8. **Status atualizado** → Database

---

## 🔧 Comandos Úteis

```bash
# Ver estrutura
tree -I 'vendor|node_modules'

# Build
make build

# Testes
make test

# Linting
make lint

# Formatar código
make format

# Limpar
make clean

# Info do projeto
make info
```

---

## ⚠️ Troubleshooting

### Erro: "Port already in use"

```bash
# Achar processo na porta
lsof -i :8080

# Matar processo
kill -9 PID
```

### Erro: Database connection

```bash
# Verificar variáveis de ambiente
cat .env

# Reset database
docker-compose down -v
docker-compose up -d
```

### Erro: Queue service

```bash
# Verificar RabbitMQ
docker-compose logs rabbitmq

# Reset
docker-compose restart rabbitmq
```

---

## 🎯 Próximos Passos

1. ✅ Arquitetura configurada
2. ✅ Banco de dados pronto
3. ✅ API funcionando
4. ✅ Workers processando
5. ⚠️ **TODO:** Implementar testes
6. ⚠️ **TODO:** Adicionar gRPC
7. ⚠️ **TODO:** Implementar cache
8. ⚠️ **TODO:** Setup CI/CD

---

**Dúvidas? Abra uma issue no GitHub!** 🙋
