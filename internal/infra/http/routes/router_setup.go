package routes

import (
	"github.com/Jeielsantosdev/Work_processamento_assincrono/internal/container"
	"github.com/Jeielsantosdev/Work_processamento_assincrono/internal/usecase"
)

// SetupRoutes configura todas as rotas da aplicação
func SetupRoutes(c *container.Container) *AppRouter {
	router := NewAppRouter()

	// Health check
	router.GET("/health", func(ctx *Context) error {
		return ctx.JSON(200, map[string]string{
			"status": "ok",
			"app":    c.Config.App.Name,
		})
	})

	// Auth routes
	router.POST("/auth/login", func(ctx *Context) error {
		return loginHandler(ctx, c)
	})

	router.POST("/auth/register", func(ctx *Context) error {
		return registerHandler(ctx, c)
	})

	// Transaction routes
	router.POST("/transactions", func(ctx *Context) error {
		return createTransactionHandler(ctx, c)
	})

	router.GET("/transactions/:id", func(ctx *Context) error {
		return getTransactionHandler(ctx, c)
	})

	router.GET("/transactions", func(ctx *Context) error {
		return listTransactionsHandler(ctx, c)
	})

	// Audit routes
	router.GET("/audit/logs", func(ctx *Context) error {
		return getAuditLogsHandler(ctx, c)
	})

	router.GET("/audit/transactions/:id", func(ctx *Context) error {
		return getTransactionAuditHandler(ctx, c)
	})

	return router
}

// Inline handlers para evitar ciclo de importação
func loginHandler(ctx *Context, c *container.Container) error {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := ctx.BindJSON(&req); err != nil {
		return ctx.JSON(400, map[string]string{"error": "Invalid request body"})
	}

	output, err := c.AuthenticateUserUC.Execute(ctx.Request.Context(), req.Email, req.Password)
	if err != nil {
		return ctx.JSON(401, map[string]string{"error": err.Error()})
	}

	return ctx.JSON(200, output)
}

func registerHandler(ctx *Context, c *container.Container) error {
	return ctx.JSON(200, map[string]string{"message": "User registration not yet implemented"})
}

func createTransactionHandler(ctx *Context, c *container.Container) error {
	var req struct {
		ClientID           string  `json:"client_id"`
		SourceAccount      string  `json:"source_account"`
		DestinationAccount string  `json:"destination_account"`
		Amount             float64 `json:"amount"`
		Currency           string  `json:"currency"`
		Description        string  `json:"description"`
	}

	if err := ctx.BindJSON(&req); err != nil {
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

	output, err := c.CreateTransactionUC.Execute(ctx.Request.Context(), input)
	if err != nil {
		return ctx.JSON(400, map[string]string{"error": err.Error()})
	}

	ctx.Writer.WriteHeader(201)
	return ctx.JSON(201, output)
}

func getTransactionHandler(ctx *Context, c *container.Container) error {
	txID := ctx.Request.URL.Query().Get("id")
	if txID == "" {
		return ctx.JSON(400, map[string]string{"error": "Transaction ID is required"})
	}

	tx, err := c.TransactionRepo.FindByID(ctx.Request.Context(), txID)
	if err != nil {
		return ctx.JSON(404, map[string]string{"error": "Transaction not found"})
	}

	return ctx.JSON(200, tx)
}

func listTransactionsHandler(ctx *Context, c *container.Container) error {
	clientID := ctx.Request.URL.Query().Get("client_id")
	if clientID == "" {
		return ctx.JSON(400, map[string]string{"error": "client_id is required"})
	}

	transactions, err := c.TransactionRepo.FindByClientID(ctx.Request.Context(), clientID)
	if err != nil {
		return ctx.JSON(500, map[string]string{"error": err.Error()})
	}

	return ctx.JSON(200, map[string]interface{}{
		"transactions": transactions,
		"count":        len(transactions),
	})
}

func getAuditLogsHandler(ctx *Context, c *container.Container) error {
	limit := 100
	offset := 0

	logs, err := c.AuditLogRepo.ListAll(ctx.Request.Context(), limit, offset)
	if err != nil {
		return ctx.JSON(500, map[string]string{"error": err.Error()})
	}

	return ctx.JSON(200, map[string]interface{}{
		"logs":  logs,
		"count": len(logs),
	})
}

func getTransactionAuditHandler(ctx *Context, c *container.Container) error {
	txID := ctx.Request.URL.Query().Get("id")
	if txID == "" {
		return ctx.JSON(400, map[string]string{"error": "Transaction ID is required"})
	}

	logs, err := c.AuditLogRepo.FindByTransactionID(ctx.Request.Context(), txID)
	if err != nil {
		return ctx.JSON(500, map[string]string{"error": err.Error()})
	}

	return ctx.JSON(200, map[string]interface{}{
		"logs":  logs,
		"count": len(logs),
	})
}
