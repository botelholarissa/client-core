package handlers

import (
	"net/http"

	"client-core/internal/models"
	"client-core/internal/service"

	"github.com/gin-gonic/gin"
)

type ClientHandler struct {
	service *service.ClientService
}

func NewClientHandler(clientService *service.ClientService) *ClientHandler {
	return &ClientHandler{service: clientService}
}

func (h *ClientHandler) Create(c *gin.Context) {
	var request models.CreateClientRequest

	err := c.ShouldBindJSON(&request)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "payload inválido",
		})

		return
	}

	mutation, err := h.service.CreateClient(request)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})

		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":         "client created successfully",
		"pipefy_mutation": mutation,
	})
}

func (h *ClientHandler) GetByEmail(c *gin.Context) {
	email := c.Param("email")

	client, err := h.service.GetClient(email)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":       client.ID,
		"nome":     client.Name,
		"email":    client.Email,
		"patrimonio": client.Assets,
		"status":   client.Status,
		"prioridade": client.Priority,
	})
}
