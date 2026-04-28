# 🔒 GUIA DE SEGURANÇA - Sistema de Auditoria Distribuída

## Índice

1. [Conceitos de Segurança](#conceitos-de-segurança)
2. [Autenticação](#autenticação)
3. [Autorização](#autorização)
4. [Criptografia](#criptografia)
5. [Validação de Input](#validação-de-input)
6. [Rate Limiting](#rate-limiting)
7. [CORS](#cors)
8. [Proteção contra Ataques](#proteção-contra-ataques)
9. [Compliance](#compliance)
10. [Security Checklist](#security-checklist)

---

## 1. CONCEITOS DE SEGURANÇA

### Princípios Core

✅ **Principle of Least Privilege**
- Usuários têm apenas as permissões necessárias
- Roles bem-definidas: Client, Auditor, Administrator

✅ **Defense in Depth**
- Múltiplas camadas de segurança
- Input validation → Authorization → Encryption

✅ **Secure by Default**
- Todas as configurações iniciam seguras
- HTTPS obrigatório em produção

✅ **Auditoria Total**
- Todas as ações registradas
- Rastreabilidade completa

---

## 2. AUTENTICAÇÃO

### JWT (JSON Web Tokens)

```go
// internal/infra/auth/jwt.go
package auth

import (
	"fmt"
	"time"
	jwt "github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID string
	Email  string
	Role   string
	jwt.RegisteredClaims
}

func GenerateToken(userID, email, role string, secret string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		UserID: userID,
		Email:  email,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	})
	
	return token.SignedString([]byte(secret))
}

func ValidateToken(tokenString, secret string) (*Claims, error) {
	claims := &Claims{}
	
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	
	if err != nil {
		return nil, err
	}
	
	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}
	
	return claims, nil
}
```

### Proteção contra Token Hijacking

```go
// Incluir fingerprint do cliente no token
type Claims struct {
	UserID      string
	Email       string
	Role        string
	Fingerprint string // Hash do User-Agent + IP
	jwt.RegisteredClaims
}

// Validar fingerprint em cada requisição
func ValidateFingerprint(token *Claims, userAgent, ip string) bool {
	expected := hashFingerprint(userAgent, ip)
	return token.Fingerprint == expected
}
```

### Password Hashing

```go
// internal/infra/crypto/bcrypt.go
package crypto

import "golang.org/x/crypto/bcrypt"

type BcryptPasswordService struct {
	cost int
}

func (s *BcryptPasswordService) HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func (s *BcryptPasswordService) VerifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword(
		[]byte(hashedPassword),
		[]byte(password),
	)
}
```

### Password Policy

```go
// Validação de senha forte
func ValidatePassword(password string) error {
	if len(password) < 12 {
		return fmt.Errorf("password must be at least 12 characters")
	}
	
	// Verificar maiúscula, minúscula, número, caractere especial
	if !regexp.MustCompile(`[A-Z]`).MatchString(password) {
		return fmt.Errorf("password must contain uppercase letter")
	}
	if !regexp.MustCompile(`[a-z]`).MatchString(password) {
		return fmt.Errorf("password must contain lowercase letter")
	}
	if !regexp.MustCompile(`[0-9]`).MatchString(password) {
		return fmt.Errorf("password must contain number")
	}
	if !regexp.MustCompile(`[!@#$%^&*]`).MatchString(password) {
		return fmt.Errorf("password must contain special character")
	}
	
	return nil
}
```

---

## 3. AUTORIZAÇÃO

### Role-Based Access Control (RBAC)

```go
// internal/domain/entities/user.go
type Role string

const (
	RoleClient        Role = "client"
	RoleAuditor       Role = "auditor"
	RoleAdministrator Role = "administrator"
)

type Permission string

const (
	PermissionCreateTransaction Permission = "create:transaction"
	PermissionViewTransaction   Permission = "view:transaction"
	PermissionViewAudit         Permission = "view:audit"
	PermissionManageUsers       Permission = "manage:users"
)

// Role permissions mapping
var rolePermissions = map[Role][]Permission{
	RoleClient: {
		PermissionCreateTransaction,
		PermissionViewTransaction,
	},
	RoleAuditor: {
		PermissionViewTransaction,
		PermissionViewAudit,
	},
	RoleAdministrator: {
		PermissionCreateTransaction,
		PermissionViewTransaction,
		PermissionViewAudit,
		PermissionManageUsers,
	},
}

func (u *User) HasPermission(permission Permission) bool {
	permissions, exists := rolePermissions[u.Role]
	if !exists {
		return false
	}
	
	for _, p := range permissions {
		if p == permission {
			return true
		}
	}
	return false
}
```

### Middleware de Autorização

```go
// internal/infra/http/middleware/auth.go
package middleware

import (
	"net/http"
	"strings"
	"fmt"
)

func AuthMiddleware(secret string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Missing authorization header", http.StatusUnauthorized)
				return
			}
			
			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || parts[0] != "Bearer" {
				http.Error(w, "Invalid authorization header", http.StatusUnauthorized)
				return
			}
			
			token := parts[1]
			claims, err := auth.ValidateToken(token, secret)
			if err != nil {
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}
			
			// Store claims in request context
			ctx := context.WithValue(r.Context(), "claims", claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func RequirePermission(permission Permission) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims := r.Context().Value("claims").(*auth.Claims)
			
			user := &entities.User{Role: entities.Role(claims.Role)}
			if !user.HasPermission(permission) {
				http.Error(w, "Insufficient permissions", http.StatusForbidden)
				return
			}
			
			next.ServeHTTP(w, r)
		})
	}
}
```

---

## 4. CRIPTOGRAFIA

### Criptografia em Trânsito (HTTPS/TLS)

```go
// cmd/api/main.go
func main() {
	// ...
	
	server := &http.Server{
		Addr:      ":8080",
		Handler:   mux,
		TLSConfig: &tls.Config{
			MinVersion: tls.VersionTLS13, // TLS 1.3 obrigatório
			CipherSuites: []uint16{
				tls.TLS_AES_256_GCM_SHA384,
				tls.TLS_CHACHA20_POLY1305_SHA256,
			},
		},
	}
	
	// Em produção, usar certificados do Let's Encrypt
	log.Fatal(server.ListenAndServeTLS("cert.pem", "key.pem"))
}
```

### Criptografia em Repouso

```go
// Dados sensíveis encriptados no banco de dados
type EncryptedField struct {
	value string
	key   []byte
}

func (e *EncryptedField) Encrypt(plaintext string) error {
	block, err := aes.NewCipher(e.key)
	if err != nil {
		return err
	}
	
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return err
	}
	
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return err
	}
	
	e.value = hex.EncodeToString(append(
		nonce,
		gcm.Seal(nil, nonce, []byte(plaintext), nil)...,
	))
	return nil
}

func (e *EncryptedField) Decrypt() (string, error) {
	data, err := hex.DecodeString(e.value)
	if err != nil {
		return "", err
	}
	
	block, err := aes.NewCipher(e.key)
	if err != nil {
		return "", err
	}
	
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	
	nonce, ciphertext := data[:gcm.NonceSize()], data[gcm.NonceSize():]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	
	return string(plaintext), err
}
```

---

## 5. VALIDAÇÃO DE INPUT

### Whitelist Validation

```go
// Validar apenas o que é esperado
func ValidateTransactionInput(input *usecase.CreateTransactionInput) error {
	// Account format: ACC + 6 números
	if !regexp.MustCompile(`^ACC\d{6}$`).MatchString(input.SourceAccount) {
		return fmt.Errorf("invalid source account format")
	}
	if !regexp.MustCompile(`^ACC\d{6}$`).MatchString(input.DestinationAccount) {
		return fmt.Errorf("invalid destination account format")
	}
	
	// Amount: 0 < amount <= 999999.99
	if input.Amount <= 0 || input.Amount > 999999.99 {
		return fmt.Errorf("invalid amount")
	}
	
	// Currency: apenas BRL, USD, EUR
	validCurrencies := map[string]bool{
		"BRL": true,
		"USD": true,
		"EUR": true,
	}
	if !validCurrencies[input.Currency] {
		return fmt.Errorf("invalid currency")
	}
	
	// IP Address
	if net.ParseIP(input.IPAddress) == nil {
		return fmt.Errorf("invalid IP address")
	}
	
	return nil
}
```

### SQL Injection Prevention

```go
// SEMPRE usar prepared statements
func (r *TransactionRepositorySQL) Save(ctx context.Context, tx *entities.Transaction) error {
	query := `
		INSERT INTO transactions (id, client_id, amount, status)
		VALUES (?, ?, ?, ?)
	`
	
	// ❌ NUNCA fazer assim:
	// query := fmt.Sprintf("INSERT INTO transactions VALUES ('%s', '%s', %.2f, '%s')", ...)
	
	// ✅ Sempre usar placeholders:
	_, err := r.db.ExecContext(ctx, query,
		tx.ID, tx.ClientID, tx.Amount, tx.Status,
	)
	
	return err
}
```

### XSS Prevention

```go
// Sanitize output
import "html"

func RenderTransaction(tx *entities.Transaction) string {
	return fmt.Sprintf(`
		<div>
			<p>ID: %s</p>
			<p>Amount: %.2f</p>
			<p>Description: %s</p>
		</div>
	`,
		html.EscapeString(tx.ID),
		tx.Amount,
		html.EscapeString(tx.Description),
	)
}
```

---

## 6. RATE LIMITING

```go
// internal/infra/http/middleware/rate_limit.go
package middleware

import (
	"net/http"
	"time"
	"golang.org/x/time/rate"
)

var limiters = make(map[string]*rate.Limiter)

func RateLimitMiddleware(requestsPerSecond int) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Usar IP do cliente como chave
			clientIP := r.RemoteAddr
			
			limiter, exists := limiters[clientIP]
			if !exists {
				limiter = rate.NewLimiter(
					rate.Limit(requestsPerSecond),
					requestsPerSecond,
				)
				limiters[clientIP] = limiter
			}
			
			if !limiter.Allow() {
				w.Header().Set("Retry-After", "60")
				http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
				return
			}
			
			next.ServeHTTP(w, r)
		})
	}
}

// Cleanup obsoleto limiteers a cada 5 minutos
func CleanupOldLimiters() {
	ticker := time.NewTicker(5 * time.Minute)
	go func() {
		for range ticker.C {
			// Remover limiters não usados
			// (implementação de limpeza)
		}
	}()
}
```

---

## 7. CORS

```go
// internal/infra/http/middleware/cors.go
package middleware

import "net/http"

func CORSMiddleware(allowedOrigins []string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")
			
			// Whitelist de origens
			allowed := false
			for _, o := range allowedOrigins {
				if origin == o {
					allowed = true
					break
				}
			}
			
			if allowed {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
				w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
				w.Header().Set("Access-Control-Max-Age", "3600")
			}
			
			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusOK)
				return
			}
			
			next.ServeHTTP(w, r)
		})
	}
}
```

---

## 8. PROTEÇÃO CONTRA ATAQUES

### CSRF Protection

```go
// CSRF token generation
func GenerateCSRFToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.StdEncoding.EncodeToString(b)
}

// Middleware
func CSRFMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" || r.Method == "PUT" || r.Method == "DELETE" {
			token := r.Header.Get("X-CSRF-Token")
			sessionToken := r.Header.Get("X-Session-Token")
			
			if token == "" || !validateCSRFToken(token, sessionToken) {
				http.Error(w, "Invalid CSRF token", http.StatusForbidden)
				return
			}
		}
		
		next.ServeHTTP(w, r)
	})
}
```

### DDOS Mitigation

```go
// Use WAF (Web Application Firewall)
// Cloudflare, AWS WAF, etc.

// Rate limiting (veja seção 6)
// IP Blocking
// Traffic Analysis
```

### Security Headers

```go
func SecurityHeadersMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Prevent clickjacking
		w.Header().Set("X-Frame-Options", "DENY")
		
		// Prevent MIME type sniffing
		w.Header().Set("X-Content-Type-Options", "nosniff")
		
		// Enable XSS protection
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		
		// Content Security Policy
		w.Header().Set("Content-Security-Policy", "default-src 'self'")
		
		// Strict Transport Security
		w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		
		// Referrer Policy
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
		
		// Feature Policy
		w.Header().Set("Permissions-Policy", "geolocation=(), camera=(), microphone=()")
		
		next.ServeHTTP(w, r)
	})
}
```

---

## 9. COMPLIANCE

### LGPD (Lei Geral de Proteção de Dados)

```go
// Right to access
func GetUserData(userID string) (*UserData, error) {
	// Retornar todos os dados do usuário
}

// Right to be forgotten
func DeleteUserData(userID string) error {
	// Deletar todos os dados pessoais do usuário
	// Manter apenas transações anônimas para auditoria
}

// Data portability
func ExportUserData(userID string) ([]byte, error) {
	// Exportar dados em formato aberto (JSON)
}
```

### PCI-DSS (Payment Card Industry)

- ✅ Não armazenar números de cartão
- ✅ Usar tokenização (Stripe, Square)
- ✅ Criptografia TLS 1.2+
- ✅ Logging completo de transações
- ✅ Testes de penetração anuais

### Auditoria Completa

```go
// Todas as ações críticas registradas
func LogAction(action string, actor *entities.User, details map[string]interface{}) {
	auditLog := &entities.AuditLog{
		ID: uuid.New().String(),
		Action: action,
		ActorID: actor.ID,
		ActorRole: string(actor.Role),
		Timestamp: time.Now(),
		Details: details,
	}
	
	// Persistir no banco
	auditLogRepository.Save(ctx, auditLog)
	
	// Enviar para SIEM (Security Information and Event Management)
	// datadog.Log(auditLog)
}
```

---

## 10. SECURITY CHECKLIST

### Desenvolvimento

- [ ] Usar HTTPS em desenvolvimento (self-signed cert ok)
- [ ] Validar todos os inputs
- [ ] Usar prepared statements
- [ ] Hash de senhas com bcrypt
- [ ] Logs de segurança habilitados
- [ ] Não commitar secrets no Git

### Deployment

- [ ] Certificado SSL válido
- [ ] HTTPS obrigatório
- [ ] HSTS header
- [ ] Firewall configurado
- [ ] Rate limiting ativo
- [ ] WAF (Web Application Firewall)
- [ ] DDoS protection

### Monitoramento

- [ ] Alertas de atividade suspeita
- [ ] Logs centralizados
- [ ] Análise de penetração mensal
- [ ] Backup de segurança testado
- [ ] Rotação de keys/secrets
- [ ] Audit log retention (7 anos mínimo)

### Compliance

- [ ] LGPD compliance verificado
- [ ] PCI-DSS compliance verificado
- [ ] SOC2 audit
- [ ] Política de privacidade publicada
- [ ] Termos de serviço atualizados
- [ ] Data Processing Agreement (DPA)

---

**Referências:**
- OWASP Top 10: https://owasp.org/Top10/
- CWE Top 25: https://cwe.mitre.org/top25/
- NIST Cybersecurity: https://www.nist.gov/cyberframework/
