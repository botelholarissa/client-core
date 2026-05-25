package main

import (
	"log"

	"client-core/internal/database"
	"client-core/internal/handlers"
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

	clientService := service.NewClientService(clientRepository)

	clientHandler := handlers.NewClientHandler(clientService)

	router := gin.Default()

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "API is running",
		})
	})

	router.POST("/clients", clientHandler.Create)

	router.Run(":8080")

}
