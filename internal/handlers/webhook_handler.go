package handlers

import (
	"net/http"

	"client-core/internal/models"
	"client-core/internal/service"

	"github.com/gin-gonic/gin"
)

type WebhookHandler struct {
	service *service.WebhookService
}

func NewWebhookHandler(s *service.WebhookService) *WebhookHandler {
	return &WebhookHandler{service: s}
}

func (h *WebhookHandler) CardUpdated(c *gin.Context) {
	var req models.PipefyWebhookRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "payload inválido"})
		return
	}

	mutation, err := h.service.ProcessWebhook(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"mutation": mutation})
}
