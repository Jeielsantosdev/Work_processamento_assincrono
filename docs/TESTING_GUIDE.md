# 🧪 GUIA DE TESTES - Sistema de Auditoria Distribuída

## Estrutura de Testes

Este guia descreve a estratégia de testes para a plataforma.

---

## 1. TESTES UNITÁRIOS

### Padrão de Teste

```go
package usecase_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/Jeielsantosdev/Work_processamento_assincrono/internal/usecase"
)

func TestCreateTransactionUC_ValidInput(t *testing.T) {
	// Arrange
	input := &usecase.CreateTransactionInput{
		ClientID: "client-123",
		SourceAccount: "ACC001",
		DestinationAccount: "ACC002",
		Amount: 100.00,
		Currency: "BRL",
		IPAddress: "192.168.1.1",
	}

	// Act
	output, err := usecase.CreateTransaction(context.Background(), input)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, output)
	assert.NotEmpty(t, output.TransactionID)
}
```

### Executar Testes
```bash
make test
```

### Cobertura de Testes
```bash
make test-coverage
```

---

## 2. TESTES DE INTEGRAÇÃO

### Setup de Teste

```go
func setupTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite", ":memory:")
	assert.NoError(t, err)
	
	// Executar migrations
	err = runMigrations(db)
	assert.NoError(t, err)
	
	return db
}
```

### Teste de Repository

```go
func TestTransactionRepository_Save(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewTransactionRepositorySQL(db)
	
	tx := &entities.Transaction{
		ID: uuid.New().String(),
		ClientID: "client-123",
		Amount: 100.00,
		Status: "PENDING",
	}
	
	err := repo.Save(context.Background(), tx)
	assert.NoError(t, err)
	
	// Verificar se foi salvo
	saved, _ := repo.FindByID(context.Background(), tx.ID)
	assert.Equal(t, tx.ID, saved.ID)
}
```

---

## 3. TESTES END-TO-END

### Teste de Fluxo Completo

```bash
#!/bin/bash

# 1. Registrar usuário
curl -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "SecurePass123!",
    "role": "client"
  }'

# 2. Login
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "SecurePass123!"
  }'

# 3. Criar transação
curl -X POST http://localhost:8080/transactions \
  -H "Authorization: Bearer {token}" \
  -H "Content-Type: application/json" \
  -d '{
    "sourceAccount": "ACC001",
    "destinationAccount": "ACC002",
    "amount": 100.00,
    "currency": "BRL",
    "description": "Payment test"
  }'

# 4. Verificar status
curl -X GET http://localhost:8080/transactions/{transactionID} \
  -H "Authorization: Bearer {token}"

# 5. Listar auditoria
curl -X GET http://localhost:8080/audit/transactions/{transactionID} \
  -H "Authorization: Bearer {token}"
```

---

## 4. TESTES DE CARGA

### Apache JMeter Test Plan

```xml
<jmeterTestPlan version="1.2">
  <hashTree>
    <TestPlan guiclass="TestPlanGui">
      <elementProp name="TestPlan.user_defined_variables" elementType="Arguments"/>
      <stringProp name="TestPlan.name">Load Test</stringProp>
    </TestPlan>
    
    <ThreadGroup guiclass="ThreadGroupGui">
      <stringProp name="ThreadGroup.num_threads">1000</stringProp>
      <stringProp name="ThreadGroup.ramp_time">60</stringProp>
      <elementProp name="ThreadGroup.main_controller" elementType="LoopController">
        <stringProp name="LoopController.loops">100</stringProp>
      </elementProp>
    </ThreadGroup>
    
    <HTTPSamplerProxy guiclass="HttpTestSampleGui">
      <stringProp name="HTTPSampler.domain">localhost</stringProp>
      <stringProp name="HTTPSampler.path">/transactions</stringProp>
      <stringProp name="HTTPSampler.method">POST</stringProp>
    </HTTPSamplerProxy>
  </hashTree>
</jmeterTestPlan>
```

### Executar com Locust
```bash
locust -f locustfile.py --host=http://localhost:8080 --users 1000 --spawn-rate 100
```

---

## 5. TESTES DE SEGURANÇA

### OWASP ZAP

```bash
zaproxy -cmd \
  -quickurl http://localhost:8080 \
  -quickout reports/security-report.html
```

### Teste de Autenticação

```bash
# ❌ Sem token
curl http://localhost:8080/transactions
# Esperado: 401 Unauthorized

# ❌ Token inválido
curl -H "Authorization: Bearer invalid" \
  http://localhost:8080/transactions
# Esperado: 401 Unauthorized

# ✅ Token válido
curl -H "Authorization: Bearer eyJ..." \
  http://localhost:8080/transactions
# Esperado: 200 OK
```

### Teste de Autorização

```bash
# ❌ Cliente tentando acessar dados de auditoria
curl -H "Authorization: Bearer {clientToken}" \
  http://localhost:8080/audit/all
# Esperado: 403 Forbidden

# ✅ Auditor acessando dados de auditoria
curl -H "Authorization: Bearer {auditorToken}" \
  http://localhost:8080/audit/all
# Esperado: 200 OK
```

---

## 6. TESTES DE PERFORMANCE

### Benchmark

```go
func BenchmarkCreateTransaction(b *testing.B) {
	container := setupContainer()
	uc := container.CreateTransactionUC
	
	for i := 0; i < b.N; i++ {
		uc.Execute(context.Background(), &usecase.CreateTransactionInput{
			ClientID: "client-123",
			SourceAccount: "ACC001",
			DestinationAccount: "ACC002",
			Amount: 100.00,
			Currency: "BRL",
			IPAddress: "192.168.1.1",
		})
	}
}
```

### Executar
```bash
go test -bench=. -benchmem ./...
```

### Métricas Esperadas
- Create Transaction: < 50ms
- Process Transaction: < 200ms
- Query Transaction: < 10ms
- Blockchain Record: < 30ms

---

## 7. PLANO DE COBERTURA DE TESTES

### Objetivos

| Componente | Meta | Atual |
|-----------|------|-------|
| **Entities** | 90% | - |
| **Use Cases** | 85% | - |
| **Repositories** | 80% | - |
| **Services** | 80% | - |
| **Total** | 85% | - |

### Checklist de Testes

- [ ] Validação de input em todos os use cases
- [ ] Permissões de acesso (RBAC)
- [ ] Criptografia de senhas
- [ ] Geração de tokens JWT
- [ ] Persistência em banco de dados
- [ ] Fila de processamento
- [ ] Blockchain recording
- [ ] Notificações
- [ ] Tratamento de erros
- [ ] Retry logic
- [ ] Timeout handling
- [ ] Concorrência

---

## 8. CI/CD TESTING

### GitHub Actions Workflow

```yaml
name: Tests

on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    
    services:
      postgres:
        image: postgres:15
        env:
          POSTGRES_PASSWORD: postgres
        options: >-
          --health-cmd pg_isready
          --health-interval 10s
          --health-timeout 5s
          --health-retries 5
    
    steps:
      - uses: actions/checkout@v3
      
      - name: Set up Go
        uses: actions/setup-go@v4
        with:
          go-version: 1.22
      
      - name: Run tests
        run: go test -v -cover ./...
      
      - name: Upload coverage
        uses: codecov/codecov-action@v3
        with:
          files: ./coverage.out
```

---

## 9. MATRIZ DE TESTES

```
┌─────────────────────┬────────┬──────┬─────┬──────────┐
│ Tipo de Teste       │ Unitário│ Integ│ E2E │ Frequência
├─────────────────────┼────────┼──────┼─────┼──────────┤
│ Entities            │   ✅    │  -   │  -  │ Por commit
│ Use Cases           │   ✅    │  ✅   │  -  │ Por commit
│ Repositories        │   ✅    │  ✅   │  -  │ Por commit
│ Services            │   ✅    │  ✅   │  -  │ Por commit
│ API Endpoints       │   -    │  ✅   │  ✅  │ Por commit
│ Security            │   -    │  -   │  ✅  │ Semanal
│ Performance         │   -    │  -   │  ✅  │ Semanal
│ Load Testing        │   -    │  -   │  ✅  │ Mensal
└─────────────────────┴────────┴──────┴─────┴──────────┘
```

---

## 10. BOAS PRÁTICAS

### ✅ Fazer
- ✅ Testes pequenos e focados
- ✅ Mocks para dependências externas
- ✅ Nomes descritivos de testes
- ✅ Testes antes do código (TDD)
- ✅ 85%+ de cobertura

### ❌ Não fazer
- ❌ Testes acoplados ao banco de dados
- ❌ Sleep/waits desnecessários
- ❌ Testes não-determinísticos
- ❌ Ignorar testes falhando
- ❌ Testes muito lentos (>1s cada)

---

**Próximas implementações:**
- Testes unitários para todas as entidades
- Suite de integração com banco PostgreSQL
- Testes de carga com simulação de 10K+ transações/seg
- Testes de segurança automatizados
