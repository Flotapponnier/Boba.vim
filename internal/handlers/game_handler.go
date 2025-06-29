package handlers

import (
	"net/http"
	"strconv"

	"boba-vim/internal/config"
	"boba-vim/internal/game"
	"boba-vim/internal/services"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type GameHandler struct {
	gameService *services.GameService
	cfg         *config.Config
}

func NewGameHandler(db *gorm.DB) *GameHandler {
	cfg := config.Load()
	return &GameHandler{
		gameService: services.NewGameService(db, cfg),
		cfg:         cfg,
	}
}

// SetUsername sets the username for the current session
func (gh *GameHandler) SetUsername(c *gin.Context) {
	var request struct {
		Username string `json:"username" binding:"required,min=2,max=50"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid username format",
		})
		return
	}

	// Get session
	session := sessions.Default(c)
	session.Set("username", request.Username)
	if err := session.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to save session",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":  true,
		"username": request.Username,
	})
}

// MovePlayer handles player movement with strict concurrency control
func (gh *GameHandler) MovePlayer(c *gin.Context) {
	var request struct {
		Direction string `json:"direction" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid request format",
		})
		return
	}

	// Get session token
	session := sessions.Default(c)
	sessionToken := session.Get("game_session_token")
	
	if sessionToken == nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "No active game session",
		})
		return
	}

	// Process move
	result, err := gh.gameService.ProcessMove(sessionToken.(string), request.Direction)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, result)
}

// GetGameState returns the current game state
func (gh *GameHandler) GetGameState(c *gin.Context) {
	session := sessions.Default(c)
	sessionToken := session.Get("game_session_token")
	if sessionToken == nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "No active game session",
		})
		return
	}

	result, err := gh.gameService.GetGameState(sessionToken.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, result)
}

// GetLeaderboard returns leaderboard data
func (gh *GameHandler) GetLeaderboard(c *gin.Context) {
	boardType := c.DefaultQuery("type", "time")
	limitStr := c.DefaultQuery("limit", "10")
	
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 || limit > 100 {
		limit = 10
	}

	result, err := gh.gameService.GetLeaderboard(boardType, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, result)
}

// GetAvailableMovements returns all available movement keys
func (gh *GameHandler) GetAvailableMovements(c *gin.Context) {
	movements := game.GetAvailableMovements()
	c.JSON(http.StatusOK, gin.H{
		"movements": movements,
		"total":     len(movements),
	})
}

// GetPlayerStats returns player statistics
func (gh *GameHandler) GetPlayerStats(c *gin.Context) {
	session := sessions.Default(c)
	username := session.Get("username")
	if username == nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "No username in session",
		})
		return
	}

	// TODO: Implement player stats retrieval
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"stats": gin.H{
			"username": username,
			"message":  "Player stats not implemented yet",
		},
	})
}

// PlayTutorial placeholder for tutorial mode
func (gh *GameHandler) PlayTutorial(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"success": false,
		"message": "Tutorial mode not implemented yet",
	})
}

// PlayOnline placeholder for online mode
func (gh *GameHandler) PlayOnline(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"success": false,
		"message": "Online mode not implemented yet",
	})
}