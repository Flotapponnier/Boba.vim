package handlers

import (
	"net/http"
	"boba-vim/internal/config"
	"boba-vim/internal/services"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type WebHandler struct {
	gameService *services.GameService
	cfg         *config.Config
}

func NewWebHandler(db *gorm.DB) *WebHandler {
	cfg := config.Load()
	return &WebHandler{
		gameService: services.NewGameService(db, cfg),
		cfg:         cfg,
	}
}

// Index serves the main page
func (wh *WebHandler) Index(c *gin.Context) {
	c.HTML(http.StatusOK, "index_go.html", gin.H{
		"title": "Boba.vim - Home",
	})
}

// PlayGame initializes a new game and serves the game page
func (wh *WebHandler) PlayGame(c *gin.Context) {
	session := sessions.Default(c)
	username := session.Get("username")
	if username == nil {
		username = "Anonymous"
	}

	// Get selected character from query parameter, default to "boba"
	selectedCharacter := c.DefaultQuery("character", "boba")

	// Create new game
	result, err := wh.gameService.CreateNewGame(username.(string), selectedCharacter)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "500_go.html", gin.H{
			"error": "Failed to initialize game: " + err.Error(),
		})
		return
	}

	if !result["success"].(bool) {
		c.HTML(http.StatusInternalServerError, "500_go.html", gin.H{
			"error": "Failed to create game session",
		})
		return
	}

	// Store session token
	session.Set("game_session_token", result["session_token"])
	session.Set("username", username)
	if err := session.Save(); err != nil {
		c.HTML(http.StatusInternalServerError, "500_go.html", gin.H{
			"error": "Failed to save session",
		})
		return
	}

	gameData := result["game_data"].(map[string]interface{})

	c.HTML(http.StatusOK, "game_go.html", gin.H{
		"title":              "Boba.vim - Game",
		"text_grid":          gameData["text_grid"],
		"game_map":           gameData["game_map"],
		"score":              gameData["score"],
		"selected_character": gameData["selected_character"],
	})
}

// NotFound serves the 404 page
func (wh *WebHandler) NotFound(c *gin.Context) {
	c.HTML(http.StatusNotFound, "404_go.html", gin.H{
		"title": "Boba.vim - Page Not Found",
		"path":  c.Request.URL.Path,
	})
}

// InternalServerError serves the 500 page
func (wh *WebHandler) InternalServerError(c *gin.Context) {
	c.HTML(http.StatusInternalServerError, "500_go.html", gin.H{
		"title": "Boba.vim - Internal Server Error",
		"error": "An unexpected error occurred. Please try again later.",
	})
}
