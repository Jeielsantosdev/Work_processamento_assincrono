package handlers

import (
	"encoding/json"

	"github.com/Jeielsantosdev/Work_processamento_assincrono/internal/container"
	"github.com/Jeielsantosdev/Work_processamento_assincrono/internal/infra/http/routes"
	"github.com/Jeielsantosdev/Work_processamento_assincrono/internal/usecase"
)

// TransactionHandler manipula requisições de transações
type TransactionHandler struct {
	container *container.Container
}

// NewTransactionHandler cria um novo manipulador de transações
func NewTransactionHandler(c *container.Container) *TransactionHandler {
	return &TransactionHandler{container: c}
}

// CreateTransaction cria uma nova transação
func (h *TransactionHandler) CreateTransaction(ctx *routes.Context) error {
	var req struct {
		ClientID           string  `json:"client_id"`
		SourceAccount      string  `json:"source_account"`
		DestinationAccount string  `json:"destination_account"`
		Amount             float64 `json:"amount"`
		Currency           string  `json:"currency"`
		Description        string  `json:"description"`
	}

	if err := json.NewDecoder(ctx.Request.Body).Decode(&req); err != nil {
		ctx.Writer.WriteHeader(400)
		return ctx.JSON(400, map[string]string{"error": "Invalid request body"})
	}

	input := &usecase.CreateTransactionInput{
		ClientID:           req.ClientID,
		SourceAccount:      req.SourceAccount,
		DestinationAccount: req.DestinationAccount,
		Amount:             req.Amount,
		Currency:           req.Currency,
		Description:        req.Description,
		IPAddress:          ctx.Request.RemoteAddr,
	}

	output, err := h.container.CreateTransactionUC.Execute(ctx.Request.Context(), input)
	if err != nil {
		ctx.Writer.WriteHeader(400)
		return ctx.JSON(400, map[string]string{"error": err.Error()})
	}

	ctx.Writer.WriteHeader(201)
	return ctx.JSON(201, output)
}

// GetTransaction obtém uma transação por ID
func (h *TransactionHandler) GetTransaction(ctx *routes.Context) error {
	txID := ctx.Params["id"]
	if txID == "" {
		return ctx.JSON(400, map[string]string{"error": "Transaction ID is required"})
	}

	tx, err := h.container.TransactionRepo.FindByID(ctx.Request.Context(), txID)
	if err != nil {
		return ctx.JSON(404, map[string]string{"error": "Transaction not found"})
	}

	return ctx.JSON(200, tx)
}

// ListTransactions lista as transações do cliente
func (h *TransactionHandler) ListTransactions(ctx *routes.Context) error {
	clientID := ctx.Request.URL.Query().Get("client_id")
	if clientID == "" {
		return ctx.JSON(400, map[string]string{"error": "client_id is required"})
	}

	// TODO: Adicionar autenticação para verificar se o cliente está listando suas próprias transações

	transactions, err := h.container.TransactionRepo.FindByClientID(ctx.Request.Context(), clientID)
	if err != nil {
		return ctx.JSON(500, map[string]string{"error": err.Error()})
	}

	return ctx.JSON(200, map[string]interface{}{
		"transactions": transactions,
		"count":        len(transactions),
	})
}
