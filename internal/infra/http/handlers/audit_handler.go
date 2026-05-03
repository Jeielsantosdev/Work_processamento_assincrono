package handlers

import (
	"github.com/Jeielsantosdev/Work_processamento_assincrono/internal/container"
	"github.com/Jeielsantosdev/Work_processamento_assincrono/internal/infra/http/routes"
)

// AuditHandler manipula requisições de auditoria
type AuditHandler struct {
	container *container.Container
}

// NewAuditHandler cria um novo manipulador de auditoria
func NewAuditHandler(c *container.Container) *AuditHandler {
	return &AuditHandler{container: c}
}

// GetAuditLogs obtém os logs de auditoria
func (h *AuditHandler) GetAuditLogs(ctx *routes.Context) error {
	// TODO: Verificar se o usuário é auditor ou administrador

	limit := 100
	offset := 0

	logs, err := h.container.AuditLogRepo.ListAll(ctx.Request.Context(), limit, offset)
	if err != nil {
		return ctx.JSON(500, map[string]string{"error": err.Error()})
	}

	return ctx.JSON(200, map[string]interface{}{
		"logs":  logs,
		"count": len(logs),
	})
}

// GetTransactionAudit obtém os logs de auditoria de uma transação
func (h *AuditHandler) GetTransactionAudit(ctx *routes.Context) error {
	txID := ctx.Params["id"]
	if txID == "" {
		return ctx.JSON(400, map[string]string{"error": "Transaction ID is required"})
	}

	logs, err := h.container.AuditLogRepo.FindByTransactionID(ctx.Request.Context(), txID)
	if err != nil {
		return ctx.JSON(500, map[string]string{"error": err.Error()})
	}

	return ctx.JSON(200, map[string]interface{}{
		"logs":  logs,
		"count": len(logs),
	})
}
