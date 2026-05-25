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
			"error": "invalid request payload",
		})

		return
	}

	err = h.service.CreateClient(request)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})

		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "client created successfully",
	})
}