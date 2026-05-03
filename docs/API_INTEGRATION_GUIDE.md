# 📡 GUIA DE INTEGRAÇÃO - Sistema de Auditoria Distribuída

## Índice

1. [REST API](#rest-api)
2. [Autenticação](#autenticação)
3. [Endpoints](#endpoints)
4. [Webhooks](#webhooks)
5. [Exemplos de Integração](#exemplos-de-integração)
6. [SDKs](#sdks)
7. [Troubleshooting](#troubleshooting)

---

## 1. REST API

### URL Base

```
Desenvolvimento: http://localhost:8080/api/v1
Staging:        https://staging-api.example.com/api/v1
Produção:       https://api.example.com/api/v1
```

### Content-Type

Todas as requisições devem enviar:
```
Content-Type: application/json
```

### Response Format

```json
{
  "success": true,
  "data": {...},
  "error": null,
  "timestamp": "2024-04-28T10:30:00Z"
}
```

### Status Codes

| Code | Significado |
|------|------------|
| **200** | OK - Requisição bem-sucedida |
| **201** | Created - Recurso criado |
| **400** | Bad Request - Input inválido |
| **401** | Unauthorized - Falta autenticação |
| **403** | Forbidden - Sem permissão |
| **404** | Not Found - Recurso não existe |
| **409** | Conflict - Conflito (ex: duplicate) |
| **429** | Too Many Requests - Rate limit |
| **500** | Server Error - Erro no servidor |
| **503** | Service Unavailable - Maintenance |

---

## 2. AUTENTICAÇÃO

### Registrar Novo Usuário

```http
POST /api/v1/auth/register
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "SecurePass123!",
  "name": "João Silva",
  "role": "client"
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "userID": "user-abc123",
    "email": "user@example.com",
    "role": "client",
    "createdAt": "2024-04-28T10:30:00Z"
  }
}
```

### Login

```http
POST /api/v1/auth/login
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "SecurePass123!"
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "expiresIn": 86400,
    "userID": "user-abc123",
    "email": "user@example.com",
    "role": "client"
  }
}
```

### Usar Token em Requisições

```http
GET /api/v1/transactions
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

---

## 3. ENDPOINTS

### Transações

#### Criar Transação

```http
POST /api/v1/transactions
Authorization: Bearer {token}
Content-Type: application/json

{
  "sourceAccount": "ACC001234",
  "destinationAccount": "ACC005678",
  "amount": 1000.50,
  "currency": "BRL",
  "description": "Payment for services",
  "metadata": {
    "invoiceID": "INV-2024-001",
    "reference": "ORD-2024-ABC"
  }
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "transactionID": "tx-abc123",
    "status": "PENDING",
    "amount": 1000.50,
    "currency": "BRL",
    "createdAt": "2024-04-28T10:30:00Z"
  }
}
```

#### Obter Transação

```http
GET /api/v1/transactions/{transactionID}
Authorization: Bearer {token}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "transactionID": "tx-abc123",
    "clientID": "client-123",
    "sourceAccount": "ACC001234",
    "destinationAccount": "ACC005678",
    "amount": 1000.50,
    "currency": "BRL",
    "status": "COMPLETED",
    "blockchainHash": "0x7f3a...",
    "retryCount": 0,
    "createdAt": "2024-04-28T10:30:00Z",
    "updatedAt": "2024-04-28T10:35:00Z",
    "completedAt": "2024-04-28T10:34:30Z"
  }
}
```

#### Listar Transações

```http
GET /api/v1/transactions?page=1&limit=20&status=COMPLETED
Authorization: Bearer {token}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "transactions": [...],
    "pagination": {
      "page": 1,
      "limit": 20,
      "total": 150,
      "pages": 8
    }
  }
}
```

### Auditoria

#### Obter Logs de Auditoria da Transação

```http
GET /api/v1/audit/transactions/{transactionID}
Authorization: Bearer {token}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "transactionID": "tx-abc123",
    "logs": [
      {
        "id": "log-1",
        "action": "TRANSACTION_CREATED",
        "actor": "client-123",
        "actorRole": "client",
        "status": "SUCCESS",
        "timestamp": "2024-04-28T10:30:00Z",
        "ipAddress": "192.168.1.100"
      },
      {
        "id": "log-2",
        "action": "TRANSACTION_PROCESSED",
        "actor": "system",
        "actorRole": "system",
        "status": "SUCCESS",
        "timestamp": "2024-04-28T10:34:30Z",
        "ipAddress": "10.0.0.1"
      }
    ]
  }
}
```

#### Listar Todos os Logs (Auditor)

```http
GET /api/v1/audit/logs?from=2024-04-01&to=2024-04-30
Authorization: Bearer {auditorToken}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "logs": [...],
    "pagination": {...}
  }
}
```

### Health Check

```http
GET /health
```

**Response:**
```json
{
  "status": "healthy",
  "database": "healthy",
  "queue": "healthy",
  "version": "1.0.0"
}
```

---

## 4. WEBHOOKS

### Registrar Webhook

```http
POST /api/v1/webhooks
Authorization: Bearer {token}
Content-Type: application/json

{
  "url": "https://yourapp.example.com/webhook",
  "events": [
    "transaction.created",
    "transaction.processing",
    "transaction.completed",
    "transaction.failed"
  ],
  "secret": "whsec_your_secret_key"
}
```

### Webhook Payload

```json
{
  "id": "evt-abc123",
  "type": "transaction.completed",
  "timestamp": "2024-04-28T10:34:30Z",
  "data": {
    "transactionID": "tx-abc123",
    "status": "COMPLETED",
    "amount": 1000.50,
    "currency": "BRL",
    "blockchainHash": "0x7f3a..."
  }
}
```

### Signature Verification

```go
import "crypto/hmac"
import "crypto/sha256"

func verifyWebhookSignature(payload []byte, signature string, secret string) bool {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write(payload)
	expectedSignature := hex.EncodeToString(h.Sum(nil))
	return hmac.Equal([]byte(signature), []byte(expectedSignature))
}
```

---

## 5. EXEMPLOS DE INTEGRAÇÃO

### cURL

```bash
#!/bin/bash

API_URL="http://localhost:8080/api/v1"
EMAIL="user@example.com"
PASSWORD="SecurePass123!"

# 1. Login
LOGIN_RESPONSE=$(curl -s -X POST "$API_URL/auth/login" \
  -H "Content-Type: application/json" \
  -d "{\"email\": \"$EMAIL\", \"password\": \"$PASSWORD\"}")

TOKEN=$(echo $LOGIN_RESPONSE | jq -r '.data.token')

# 2. Criar transação
curl -X POST "$API_URL/transactions" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "sourceAccount": "ACC001234",
    "destinationAccount": "ACC005678",
    "amount": 100.00,
    "currency": "BRL",
    "description": "Test payment"
  }'

# 3. Obter transação
curl -X GET "$API_URL/transactions/tx-abc123" \
  -H "Authorization: Bearer $TOKEN"
```

### Python

```python
import requests
import json

class AuditoriaClient:
    def __init__(self, base_url, email, password):
        self.base_url = base_url
        self.session = requests.Session()
        self.authenticate(email, password)
    
    def authenticate(self, email, password):
        response = self.session.post(
            f"{self.base_url}/auth/login",
            json={"email": email, "password": password}
        )
        self.token = response.json()["data"]["token"]
        self.session.headers.update({
            "Authorization": f"Bearer {self.token}"
        })
    
    def create_transaction(self, source, destination, amount, currency):
        response = self.session.post(
            f"{self.base_url}/transactions",
            json={
                "sourceAccount": source,
                "destinationAccount": destination,
                "amount": amount,
                "currency": currency
            }
        )
        return response.json()["data"]
    
    def get_transaction(self, transaction_id):
        response = self.session.get(
            f"{self.base_url}/transactions/{transaction_id}"
        )
        return response.json()["data"]

# Uso
client = AuditoriaClient(
    "http://localhost:8080/api/v1",
    "user@example.com",
    "SecurePass123!"
)

tx = client.create_transaction(
    "ACC001234",
    "ACC005678",
    100.00,
    "BRL"
)

print(f"Transaction ID: {tx['transactionID']}")
```

### JavaScript/Node.js

```javascript
const axios = require('axios');

class AuditoriaClient {
  constructor(baseUrl, email, password) {
    this.baseUrl = baseUrl;
    this.client = axios.create();
    this.authenticate(email, password);
  }
  
  async authenticate(email, password) {
    const response = await this.client.post(
      `${this.baseUrl}/auth/login`,
      { email, password }
    );
    
    this.token = response.data.data.token;
    this.client.defaults.headers.common['Authorization'] = `Bearer ${this.token}`;
  }
  
  async createTransaction(source, destination, amount, currency) {
    const response = await this.client.post(
      `${this.baseUrl}/transactions`,
      {
        sourceAccount: source,
        destinationAccount: destination,
        amount,
        currency
      }
    );
    
    return response.data.data;
  }
  
  async getTransaction(transactionId) {
    const response = await this.client.get(
      `${this.baseUrl}/transactions/${transactionId}`
    );
    
    return response.data.data;
  }
}

// Uso
(async () => {
  const client = new AuditoriaClient(
    'http://localhost:8080/api/v1',
    'user@example.com',
    'SecurePass123!'
  );
  
  const tx = await client.createTransaction(
    'ACC001234',
    'ACC005678',
    100.00,
    'BRL'
  );
  
  console.log(`Transaction ID: ${tx.transactionID}`);
})();
```

### Go

```go
package main

import (
	"fmt"
	"github.com/your-org/auditoria-go-sdk/client"
)

func main() {
	c := client.New("http://localhost:8080/api/v1")
	
	// Authenticate
	token, err := c.Authenticate("user@example.com", "SecurePass123!")
	if err != nil {
		panic(err)
	}
	c.SetToken(token)
	
	// Create transaction
	tx, err := c.CreateTransaction(
		client.CreateTransactionInput{
			SourceAccount:      "ACC001234",
			DestinationAccount: "ACC005678",
			Amount:             100.00,
			Currency:           "BRL",
		},
	)
	if err != nil {
		panic(err)
	}
	
	fmt.Printf("Transaction ID: %s\n", tx.TransactionID)
}
```

---

## 6. SDKs

### Disponíveis

- ✅ **Go** (`auditoria-go-sdk`)
- 🚧 **Python** (`auditoria-python-sdk`) - Em desenvolvimento
- 🚧 **Node.js** (`auditoria-node-sdk`) - Em desenvolvimento
- 🚧 **Java** (`auditoria-java-sdk`) - Em desenvolvimento

### Instalação

```bash
# Go
go get github.com/Jeielsantosdev/auditoria-go-sdk

# Python
pip install auditoria-python-sdk

# Node.js
npm install auditoria-node-sdk
```

---

## 7. TROUBLESHOOTING

### Erro: 401 Unauthorized

```
Cause: Token ausente, inválido ou expirado
Solution: 
1. Verify token in Authorization header
2. Check token expiration (24 hours default)
3. Re-authenticate to get new token
```

### Erro: 403 Forbidden

```
Cause: Usuário não tem permissão
Solution:
1. Verify user role
2. Check required permissions for endpoint
3. Contact administrator
```

### Erro: 429 Too Many Requests

```
Cause: Rate limit excedido
Solution:
1. Wait 60 seconds before retrying
2. Check rate limit headers
3. Implement exponential backoff
```

### Erro: 500 Internal Server Error

```
Cause: Erro no servidor
Solution:
1. Check server logs
2. Retry after 30 seconds
3. Contact support
```

---

**Documentação Completa:**
Https://docs.example.com/api

**Status da API:**
https://status.example.com
