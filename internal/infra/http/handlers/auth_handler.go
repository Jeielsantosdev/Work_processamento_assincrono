package handlers

import (
	"encoding/json"

	"github.com/Jeielsantosdev/Work_processamento_assincrono/internal/container"
	"github.com/Jeielsantosdev/Work_processamento_assincrono/internal/infra/http/routes"
)

// AuthHandler manipula requisições de autenticação
type AuthHandler struct {
	container *container.Container
}

// NewAuthHandler cria um novo manipulador de autenticação
func NewAuthHandler(c *container.Container) *AuthHandler {
	return &AuthHandler{container: c}
}

// Login autentica um usuário
func (h *AuthHandler) Login(ctx *routes.Context) error {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(ctx.Request.Body).Decode(&req); err != nil {
		ctx.Writer.WriteHeader(400)
		return err
	}

	output, err := h.container.AuthenticateUserUC.Execute(ctx.Request.Context(), req.Email, req.Password)
	if err != nil {
		ctx.Writer.WriteHeader(401)
		return ctx.JSON(401, map[string]string{"error": err.Error()})
	}

	return ctx.JSON(200, output)
}

// Register registra um novo usuário
func (h *AuthHandler) Register(ctx *routes.Context) error {
	var req struct {
		Email    string `json:"email"`
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(ctx.Request.Body).Decode(&req); err != nil {
		ctx.Writer.WriteHeader(400)
		return err
	}

	// TODO: Implementar use case de registro
	return ctx.JSON(200, map[string]string{"message": "User registered successfully"})
}
