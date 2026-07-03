package main

import (
	"fmt"
	"log"

	"github.com/Sheikh-Fahad-Ahmed/Team-Pulse-Metrics-Engine/internal/database"
<<<<<<< HEAD
	"github.com/Sheikh-Fahad-Ahmed/Team-Pulse-Metrics-Engine/internal/webhook/handlers"
=======
	"github.com/Sheikh-Fahad-Ahmed/Team-Pulse-Metrics-Engine/internal/handlers"
>>>>>>> main

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	database.ConnectDB()
	err := godotenv.Load()
	if err != nil {
		fmt.Println(err)
	}

	r := gin.Default()

	r.POST("/api/v1/webhook/github", handlers.HandleWebhook)

	if err := r.Run(); err != nil {
		log.Fatalf("server failed to start: %v", err)
	}
}
