package main

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"pulsepoll/backend/config"
	"pulsepoll/backend/routes"
)

func main() {
	// Load environment variables from .env
	if err := godotenv.Load(); err != nil {
		panic("Error loading .env file")
	}

	// Connect to MongoDB
	if err := config.ConnectDatabase(); err != nil {
		panic(err)
	}

	// Create Gin router
	router := gin.Default()

	// Allow React frontend to access the Go API
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	})

	// Home route
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "PulsePoll backend is running!",
		})
	})

	// Health check route
	router.GET("/api/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"message": "PulsePoll API is healthy",
		})
	})

	// Poll routes
	router.POST("/api/polls", routes.CreatePoll)
	router.GET("/api/polls", routes.GetPolls)

	// Vote route
	router.POST("/api/polls/:id/vote/:optionId", routes.VotePoll)

	// Results route
	router.GET("/api/polls/:id/results", routes.GetPollResults)

	// Get port from .env
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Start server
	if err := router.Run(":" + port); err != nil {
		panic(err)
	}
}