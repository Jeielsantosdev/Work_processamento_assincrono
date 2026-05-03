package auth

import (
	"fmt"
	"time"

	"github.com/Jeielsantosdev/Work_processamento_assincrono/internal/domain/entities"
	"github.com/Jeielsantosdev/Work_processamento_assincrono/internal/domain/interfaces"
	"github.com/golang-jwt/jwt/v5"
)

// JWTAuthService implementa o serviço de autenticação com JWT
type JWTAuthService struct {
	secretKey string
	expiresIn time.Duration
}

// NewJWTAuthService cria uma nova instância do serviço JWT
func NewJWTAuthService(secretKey string, expiresIn time.Duration) *JWTAuthService {
	if expiresIn == 0 {
		expiresIn = 24 * time.Hour
	}
	return &JWTAuthService{
		secretKey: secretKey,
		expiresIn: expiresIn,
	}
}

// GenerateToken gera um token JWT
func (s *JWTAuthService) GenerateToken(user *entities.User) (string, error) {
	now := time.Now()
	expiresAt := now.Add(s.expiresIn)

	claims := jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"role":    string(user.Role),
		"iat":     now.Unix(),
		"exp":     expiresAt.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(s.secretKey))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return signedToken, nil
}

// ValidateToken valida um token JWT
func (s *JWTAuthService) ValidateToken(tokenString string) (*interfaces.TokenClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Verificar se o método de assinatura é HS256
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.secretKey), nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	return &interfaces.TokenClaims{
		UserID:    claims["user_id"].(string),
		Email:     claims["email"].(string),
		Role:      entities.UserRole(claims["role"].(string)),
		IssuedAt:  int64(claims["iat"].(float64)),
		ExpiresAt: int64(claims["exp"].(float64)),
	}, nil
}
