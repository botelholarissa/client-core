package main

import (
	"log"

	"client-core/internal/database"
	"client-core/internal/handlers"
	"client-core/internal/pipefy"
	"client-core/internal/repository"
	"client-core/internal/service"

	"github.com/gin-gonic/gin"
)

func main() {
	db, err := database.Connect()

	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	clientRepository := repository.NewClientRepository(db)

	webhookRepository := repository.NewWebhookRepository(db)

	pipefyClient := pipefy.NewPipefyClient()
	clientService := service.NewClientService(clientRepository, pipefyClient, 34)

	clientHandler := handlers.NewClientHandler(clientService)

	webhookService := service.NewWebhookService(clientRepository, webhookRepository, pipefyClient)
	webhookHandler := handlers.NewWebhookHandler(webhookService)

	router := gin.Default()

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "API is running",
		})
	})

	router.POST("/clientes", clientHandler.Create)
	router.POST("/webhooks/pipefy/card-updated", webhookHandler.CardUpdated)

	router.Run(":8080")

}
