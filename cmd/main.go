package main

import (
	"log"
	"os"
	"strconv"

	"client-core/internal/database"
	"client-core/internal/handlers"
	"client-core/internal/pipefy"
	"client-core/internal/repository"
	"client-core/internal/service"

	"github.com/joho/godotenv"
	"github.com/gin-gonic/gin"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db, err := database.Connect()

	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	clientRepository := repository.NewClientRepository(db)

	webhookRepository := repository.NewWebhookRepository(db)

	pipefyClient := pipefy.NewPipefyClient()

	pipeIDValue := os.Getenv("PIPEFY_ID")

	if pipeIDValue == "" {
		log.Fatal("PIPEFY_ID is required")
	}

	pipeID, err := strconv.ParseInt(pipeIDValue, 10, 64)

	if err != nil {
		log.Fatal("invalid PIPEFY_ID")
	}

	clientService := service.NewClientService(clientRepository, pipefyClient, pipeID)

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
	router.GET("/clientes/:email", clientHandler.GetByEmail)
	router.POST("/webhooks/pipefy/card-updated", webhookHandler.CardUpdated)

	router.Run(":8080")

}
