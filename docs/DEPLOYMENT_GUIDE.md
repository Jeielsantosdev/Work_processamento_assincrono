# 🚀 GUIA DE DEPLOYMENT - Sistema de Auditoria Distribuída

## Índice

1. [Deployment Local](#deployment-local)
2. [Deployment em Staging](#deployment-em-staging)
3. [Deployment em Produção](#deployment-em-produção)
4. [Monitoramento](#monitoramento)
5. [Rollback](#rollback)
6. [Scaling](#scaling)
7. [Disaster Recovery](#disaster-recovery)

---

## 1. DEPLOYMENT LOCAL

### Pré-requisitos

```bash
# ✅ Instalar dependências
brew install go docker docker-compose postgresql

# ✅ Verificar versões
go version     # 1.22+
docker --version
docker-compose --version
```

### Start Rápido

```bash
# 1. Clone e configure
git clone https://github.com/Jeielsantosdev/Work_processamento_assincrono.git
cd Work_processamento_assincrono
cp .env.example .env

# 2. Inicie os containers
docker-compose up -d

# 3. Execute migrations
go run cmd/api/main.go --migrate

# 4. API disponível em http://localhost:8080
curl http://localhost:8080/health
```

### Desenvolvimento

```bash
# Terminal 1: API Server
go run cmd/api/main.go

# Terminal 2: Worker Service
go run cmd/worker/worker.go

# Terminal 3: Logs
tail -f logs/app.log
```

---

## 2. DEPLOYMENT EM STAGING

### Ambiente Staging

```yaml
# .env.staging
APP_ENV=staging
APP_PORT=8080
DB_HOST=postgres-staging.internal
DB_NAME=auditoria_staging
LOG_LEVEL=INFO
JWT_SECRET=staging-secret-key-change-in-prod
QUEUE_TYPE=inmemory
```

### Deploy via Docker

```bash
# 1. Build images
docker build -t auditoria:staging --target api .

# 2. Tag
docker tag auditoria:staging \
  registry.example.com/auditoria:staging-$(git rev-parse --short HEAD)

# 3. Push
docker push registry.example.com/auditoria:staging-$(git rev-parse --short HEAD)

# 4. Deploy
kubectl apply -f k8s/staging/deployment.yaml
```

### Kubernetes Deployment (staging)

```yaml
# k8s/staging/deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: auditoria-api-staging
  namespace: auditoria-staging
spec:
  replicas: 2
  selector:
    matchLabels:
      app: auditoria-api
      env: staging
  template:
    metadata:
      labels:
        app: auditoria-api
        env: staging
    spec:
      containers:
      - name: api
        image: registry.example.com/auditoria:staging-abc123
        ports:
        - containerPort: 8080
        env:
        - name: APP_ENV
          value: staging
        - name: DB_HOST
          valueFrom:
            secretKeyRef:
              name: db-credentials
              key: host
        livenessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 30
          periodSeconds: 10
        resources:
          requests:
            memory: "256Mi"
            cpu: "250m"
          limits:
            memory: "512Mi"
            cpu: "500m"
```

### Service (staging)

```yaml
# k8s/staging/service.yaml
apiVersion: v1
kind: Service
metadata:
  name: auditoria-api-staging
  namespace: auditoria-staging
spec:
  type: LoadBalancer
  ports:
  - port: 80
    targetPort: 8080
  selector:
    app: auditoria-api
    env: staging
```

### Smoke Tests

```bash
#!/bin/bash

STAGING_URL="https://staging-api.example.com"

# Health check
curl -f $STAGING_URL/health || exit 1

# API readiness
curl -f $STAGING_URL/api/status || exit 1

echo "✅ Staging deployment successful"
```

---

## 3. DEPLOYMENT EM PRODUÇÃO

### Pré-requisitos Produção

- ✅ Certificado SSL/TLS válido
- ✅ CDN com DDoS protection
- ✅ Backup automatizado
- ✅ Monitoramento 24/7
- ✅ Alertas configurados
- ✅ Plano de disaster recovery
- ✅ Aprovação de segurança

### Variáveis Produção

```yaml
# .env.production
APP_ENV=production
APP_PORT=8080
APP_NAME=Auditoria-Distribuida
LOG_LEVEL=WARN

# Database
DB_HOST=prod-db-cluster.rds.amazonaws.com
DB_PORT=5432
DB_USER=${DB_USER}  # From AWS Secrets Manager
DB_PASSWORD=${DB_PASSWORD}  # From AWS Secrets Manager
DB_NAME=auditoria_prod
DB_SSL_MODE=require
DB_MAX_CONNECTIONS=100

# JWT
JWT_SECRET=${JWT_SECRET}  # From AWS Secrets Manager
JWT_EXPIRATION=86400

# Queue
QUEUE_TYPE=rabbitmq
QUEUE_HOST=rabbitmq-prod.internal
QUEUE_PORT=5672

# Blockchain
BLOCKCHAIN_TYPE=ethereum
BLOCKCHAIN_NETWORK=mainnet
BLOCKCHAIN_CONTRACT=${CONTRACT_ADDRESS}  # From AWS Secrets Manager

# Monitoring
DD_ENV=production
DD_SERVICE=auditoria-api
DD_TRACE_ENABLED=true
```

### Production Deployment

```bash
#!/bin/bash

set -e

VERSION=$(git describe --tags)
REGION="us-east-1"

echo "🚀 Deploying Auditoria v$VERSION to Production"

# 1. Build
docker build -t auditoria:$VERSION --target api .
docker tag auditoria:$VERSION registry.example.com/auditoria:$VERSION

# 2. Security scan
trivy image registry.example.com/auditoria:$VERSION

# 3. Push
aws ecr get-login-password --region $REGION | docker login --username AWS --password-stdin registry.example.com
docker push registry.example.com/auditoria:$VERSION

# 4. Deploy via Blue-Green
./scripts/blue-green-deploy.sh $VERSION

# 5. Smoke tests
./scripts/smoke-tests.sh production

# 6. Monitor
./scripts/monitor-deployment.sh 300  # 5 minutes

echo "✅ Production deployment complete"
```

### Blue-Green Deployment

```bash
#!/bin/bash
# scripts/blue-green-deploy.sh

VERSION=$1

# 1. Deploy para "Green" (novo ambiente)
kubectl set image deployment/auditoria-api-green \
  auditoria=registry.example.com/auditoria:$VERSION \
  -n production

# 2. Esperar por rollout
kubectl rollout status deployment/auditoria-api-green -n production

# 3. Testes de smoke
if ./scripts/smoke-tests.sh green; then
  # 4. Switch traffic Blue -> Green
  kubectl patch service auditoria-api \
    -p '{"spec":{"selector":{"deployment":"green"}}}' \
    -n production
  
  echo "✅ Traffic switched to Green deployment"
else
  echo "❌ Smoke tests failed, keeping Blue active"
  exit 1
fi
```

### Health Checks

```go
// internal/health/health.go
package health

import (
	"net/http"
	"database/sql"
	"context"
	"encoding/json"
)

type HealthResponse struct {
	Status   string `json:"status"`
	Database string `json:"database"`
	Queue    string `json:"queue"`
	Version  string `json:"version"`
}

func HealthCheck(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		health := HealthResponse{Status: "healthy"}
		
		// Check database
		err := db.PingContext(context.Background())
		if err != nil {
			health.Database = "unhealthy"
		} else {
			health.Database = "healthy"
		}
		
		// Check queue
		// ... queue check
		
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(health)
	}
}
```

---

## 4. MONITORAMENTO

### Prometheus Metrics

```go
// internal/metrics/metrics.go
package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	TransactionsCreated = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "transactions_created_total",
			Help: "Total transactions created",
		},
		[]string{"status"},
	)
	
	ProcessingDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "transaction_processing_seconds",
			Help: "Transaction processing duration",
			Buckets: []float64{.01, .05, .1, .5, 1, 2, 5},
		},
		[]string{"operation"},
	)
)

func init() {
	prometheus.MustRegister(TransactionsCreated, ProcessingDuration)
}
```

### Datadog Dashboard

```json
{
  "title": "Auditoria API - Produção",
  "widgets": [
    {
      "definition": {
        "type": "timeseries",
        "requests": [
          {
            "q": "avg:auditoria.transaction.latency{env:prod}"
          }
        ],
        "title": "Transaction Latency"
      }
    },
    {
      "definition": {
        "type": "gauge",
        "requests": [
          {
            "q": "avg:system.cpu{host:prod-api-*}"
          }
        ],
        "title": "CPU Usage"
      }
    }
  ]
}
```

### Alertas

```yaml
# alerts.yaml
groups:
  - name: auditoria
    rules:
      - alert: HighLatency
        expr: histogram_quantile(0.95, transaction_processing_seconds) > 1
        for: 5m
        annotations:
          summary: "High transaction latency detected"
      
      - alert: DatabaseDown
        expr: up{job="postgres"} == 0
        for: 1m
        annotations:
          summary: "Database is down"
```

---

## 5. ROLLBACK

### Rollback Automático

```bash
#!/bin/bash
# scripts/rollback.sh

PREVIOUS_VERSION=$(git describe --tags --abbrev=0 HEAD^)

echo "🔄 Rolling back to $PREVIOUS_VERSION"

# 1. Get previous image
docker pull registry.example.com/auditoria:$PREVIOUS_VERSION

# 2. Switch deployment
kubectl set image deployment/auditoria-api \
  auditoria=registry.example.com/auditoria:$PREVIOUS_VERSION \
  -n production

# 3. Monitor
kubectl rollout status deployment/auditoria-api -n production

# 4. Verify
./scripts/smoke-tests.sh production

echo "✅ Rollback complete"
```

### Manual Rollback

```bash
kubectl rollout undo deployment/auditoria-api -n production
kubectl rollout status deployment/auditoria-api -n production
```

---

## 6. SCALING

### Horizontal Scaling (HPA)

```yaml
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: auditoria-api-hpa
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: auditoria-api
  minReplicas: 3
  maxReplicas: 100
  metrics:
  - type: Resource
    resource:
      name: cpu
      target:
        type: Utilization
        averageUtilization: 70
  - type: Resource
    resource:
      name: memory
      target:
        type: Utilization
        averageUtilization: 80
```

### Vertical Scaling

```yaml
# requests/limits
resources:
  requests:
    memory: "512Mi"
    cpu: "500m"
  limits:
    memory: "2Gi"
    cpu: "2000m"
```

### Database Scaling

```bash
# RDS Read Replica
aws rds create-db-instance-read-replica \
  --db-instance-identifier auditoria-read-replica \
  --source-db-identifier auditoria-prod
```

---

## 7. DISASTER RECOVERY

### Backup Strategy

```bash
#!/bin/bash
# scripts/backup.sh

BACKUP_DATE=$(date +%Y%m%d_%H%M%S)
BACKUP_PATH="s3://auditoria-backups/prod/db-$BACKUP_DATE.sql"

# 1. Dump database
pg_dump \
  -h $DB_HOST \
  -U $DB_USER \
  -d $DB_NAME \
  | gzip \
  | aws s3 cp - $BACKUP_PATH

# 2. Backup blockchain state
./scripts/backup-blockchain.sh

echo "✅ Backup completed: $BACKUP_PATH"
```

### RTO/RPO

| Cenário | RTO | RPO |
|---------|-----|-----|
| App crash | 5 min | 0 |
| Single node failure | 2 min | 0 |
| Region failure | 30 min | 5 min |
| Data corruption | 1 hour | 24 hours |

### Disaster Recovery Plan

```bash
#!/bin/bash
# scripts/disaster-recovery.sh

echo "🚨 Starting Disaster Recovery"

# 1. Restore database from backup
aws s3 cp s3://auditoria-backups/prod/db-latest.sql.gz - | gunzip | psql

# 2. Restore blockchain state
./scripts/restore-blockchain.sh

# 3. Restart services
kubectl rollout restart deployment/auditoria-api -n production
kubectl rollout restart deployment/auditoria-worker -n production

# 4. Verify
./scripts/smoke-tests.sh production

echo "✅ Disaster Recovery complete"
```

---

## Checklist de Deployment

- [ ] Testes passando (cobertura > 85%)
- [ ] Security scan sem vulnerabilidades críticas
- [ ] Staging deployment bem-sucedido
- [ ] Smoke tests passando
- [ ] Monitores configurados
- [ ] Alertas testados
- [ ] Rollback plan documentado
- [ ] Backup recente disponível
- [ ] Aprovação de security team
- [ ] Aprovação de product team

---

**Próximas implementações:**
- Terraform para infrastructure as code
- ArgoCD para GitOps
- Vault para secret management
- Detailed runbooks para cada cenário
