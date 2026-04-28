package crypto

import (
	"golang.org/x/crypto/bcrypt"
)

// BcryptPasswordService implementa o serviço de criptografia de senha
type BcryptPasswordService struct {
	cost int
}

// NewBcryptPasswordService cria uma nova instância do serviço de senha
func NewBcryptPasswordService(cost int) *BcryptPasswordService {
	if cost == 0 {
		cost = bcrypt.DefaultCost
	}
	return &BcryptPasswordService{cost: cost}
}

// HashPassword criptografa a senha
func (s *BcryptPasswordService) HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), s.cost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// VerifyPassword verifica a senha
func (s *BcryptPasswordService) VerifyPassword(hashedPassword, plainPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
	return err == nil
}
