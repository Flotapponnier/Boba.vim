package main

import (
	"log"
	"net/http"

	"boba-vim/internal/config"
	"boba-vim/internal/database"
	"boba-vim/internal/handlers"
	"boba-vim/internal/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize database
	db, err := database.Initialize(cfg.DatabaseURL)
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}

	// Initialize Gin router
	router := gin.Default()

	// Add middleware
	router.Use(middleware.CORS())
	router.Use(middleware.Sessions(cfg.SessionSecret))
	router.Use(middleware.Logger())

	// Serve static files
	router.Static("/static", "./static")
	router.LoadHTMLGlob("templates/*_go.html")

	// Initialize handlers
	gameHandler := handlers.NewGameHandler(db)
	webHandler := handlers.NewWebHandler(db)

	// Web routes
	router.GET("/", webHandler.Index)
	router.GET("/api/play", webHandler.PlayGame)

	// API routes
	api := router.Group("/api")
	{
		api.POST("/set-username", gameHandler.SetUsername)
		api.POST("/move", gameHandler.MovePlayer)
		api.GET("/game-state", gameHandler.GetGameState)
		api.GET("/leaderboard", gameHandler.GetLeaderboard)
		api.GET("/movements", gameHandler.GetAvailableMovements)
		api.GET("/player-stats", gameHandler.GetPlayerStats)
		api.POST("/playtutorial", gameHandler.PlayTutorial)
		api.POST("/playonline", gameHandler.PlayOnline)
	}

	// Start server
	log.Printf("Server starting on port %s", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, router))
}