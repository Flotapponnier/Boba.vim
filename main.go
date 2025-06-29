package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
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
	router.Use(middleware.ErrorHandler()) // Add error handling middleware

	// Serve static files
	router.Static("/static", "./static")
	router.LoadHTMLGlob("templates/*_go.html")

	// Initialize handlers
	gameHandler := handlers.NewGameHandler(db)
	webHandler := handlers.NewWebHandler(db)

	// Web routes
	router.GET("/", webHandler.Index)
	router.GET("/api/play", webHandler.PlayGame)

	// Test routes (remove in production)
	router.GET("/test-404", func(c *gin.Context) {
		c.Redirect(http.StatusFound, "/nonexistent-page")
	})
	router.GET("/test-panic", func(c *gin.Context) {
		panic("This is a test panic!")
	})
	router.GET("/test-500-page", webHandler.InternalServerError)

	// API routes
	api := router.Group("/api")
	{
		api.POST("/set-username", gameHandler.SetUsername)
		api.POST("/move", gameHandler.MovePlayer)
		api.GET("/game-state", gameHandler.GetGameState)
		api.GET("/leaderboard", gameHandler.GetLeaderboard)
		api.GET("/movements", gameHandler.GetAvailableMovements)
		api.GET("/player-stats", gameHandler.GetPlayerStats)
		api.POST("/playonline", gameHandler.PlayOnline)
		
		// Test error pages
		api.GET("/test-500", webHandler.InternalServerError)
		api.GET("/test-panic", func(c *gin.Context) {
			panic("API test panic!")
		})
	}

	// 404 handler - must be last route
	router.NoRoute(webHandler.NotFound)

	// Create HTTP server
	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: router,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Server starting on port %s", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	// Catch more signals including terminal closure
	signal.Notify(quit, 
		syscall.SIGINT,  // Ctrl+C
		syscall.SIGTERM, // kill command
		syscall.SIGHUP,  // Terminal hangup (closing terminal)
		syscall.SIGQUIT, // Ctrl+\
	)
	<-quit
	log.Println("Shutting down server...")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exited")
}
