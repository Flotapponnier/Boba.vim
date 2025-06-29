package handlers

import (
	"net/http"
	"strings"
	"golang.org/x/crypto/bcrypt"
	
	"boba-vim/internal/models"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AuthHandler struct {
	db *gorm.DB
}

func NewAuthHandler(db *gorm.DB) *AuthHandler {
	return &AuthHandler{db: db}
}

type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=20"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Register handles user registration
func (ah *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		// Provide user-friendly error messages
		errorMsg := "Please check your input"
		if strings.Contains(err.Error(), "min") && strings.Contains(err.Error(), "Username") {
			errorMsg = "Username must be at least 3 characters long"
		} else if strings.Contains(err.Error(), "min") && strings.Contains(err.Error(), "Password") {
			errorMsg = "Password must be at least 8 characters long"
		} else if strings.Contains(err.Error(), "max") && strings.Contains(err.Error(), "Username") {
			errorMsg = "Username must be less than 20 characters"
		} else if strings.Contains(err.Error(), "email") {
			errorMsg = "Please enter a valid email address"
		} else if strings.Contains(err.Error(), "required") {
			errorMsg = "All fields are required"
		}
		
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   errorMsg,
		})
		return
	}

	// Clean username and email
	req.Username = strings.TrimSpace(req.Username)
	req.Email = strings.TrimSpace(strings.ToLower(req.Email))

	// Check if user already exists
	var existingPlayer models.Player
	if err := ah.db.Where("username = ? OR email = ?", req.Username, req.Email).First(&existingPlayer).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{
			"success": false,
			"error":   "Username or email already exists",
		})
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to process password",
		})
		return
	}

	// Create new player
	player := models.Player{
		Username:     req.Username,
		Email:        req.Email,
		Password:     string(hashedPassword),
		IsRegistered: true,
	}

	if err := ah.db.Create(&player).Error; err != nil {
		errorMsg := "Failed to create account"
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			if strings.Contains(err.Error(), "username") {
				errorMsg = "Username already exists. Please choose a different one."
			} else if strings.Contains(err.Error(), "email") {
				errorMsg = "Email already registered. Please use a different email or try logging in."
			} else {
				errorMsg = "Username or email already exists"
			}
		}
		
		c.JSON(http.StatusConflict, gin.H{
			"success": false,
			"error":   errorMsg,
		})
		return
	}

	// Log the user in immediately after registration
	session := sessions.Default(c)
	session.Set("user_id", player.ID)
	session.Set("username", player.Username)
	session.Set("is_registered", true)
	if err := session.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to save session",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Account created successfully",
		"user": gin.H{
			"id":            player.ID,
			"username":      player.Username,
			"email":         player.Email,
			"is_registered": player.IsRegistered,
		},
	})
}

// Login handles user login
func (ah *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid input: " + err.Error(),
		})
		return
	}

	// Find user by username or email
	var player models.Player
	if err := ah.db.Where("username = ? OR email = ?", req.Username, req.Username).First(&player).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "Invalid credentials",
		})
		return
	}

	// Check password
	if err := bcrypt.CompareHashAndPassword([]byte(player.Password), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "Invalid credentials",
		})
		return
	}

	// Create session
	session := sessions.Default(c)
	session.Set("user_id", player.ID)
	session.Set("username", player.Username)
	session.Set("is_registered", player.IsRegistered)
	if err := session.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to save session",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Login successful",
		"user": gin.H{
			"id":            player.ID,
			"username":      player.Username,
			"email":         player.Email,
			"is_registered": player.IsRegistered,
		},
	})
}

// Logout handles user logout
func (ah *AuthHandler) Logout(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	if err := session.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to logout",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Logged out successfully",
	})
}

// GetCurrentUser returns current user info
func (ah *AuthHandler) GetCurrentUser(c *gin.Context) {
	session := sessions.Default(c)
	userID := session.Get("user_id")
	username := session.Get("username")
	isRegistered := session.Get("is_registered")

	if userID == nil {
		c.JSON(http.StatusOK, gin.H{
			"success":       true,
			"authenticated": false,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":       true,
		"authenticated": true,
		"user": gin.H{
			"id":            userID,
			"username":      username,
			"is_registered": isRegistered,
		},
	})
}