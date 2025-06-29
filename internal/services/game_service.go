package services

import (
	"errors"
	"time"

	"boba-vim/internal/config"
	"boba-vim/internal/game"
	"boba-vim/internal/models"

	"gorm.io/gorm"
)

type GameService struct {
	db  *gorm.DB
	cfg *config.Config
}

func NewGameService(db *gorm.DB, cfg *config.Config) *GameService {
	return &GameService{
		db:  db,
		cfg: cfg,
	}
}

// CreateNewGame creates a new secure game session
func (gs *GameService) CreateNewGame(username string) (map[string]interface{}, error) {
	// Get or create player
	var player models.Player
	result := gs.db.Where("username = ?", username).First(&player)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			player = models.Player{
				Username:    username,
				TotalGames:  0,
				BestScore:   0,
				TotalPearls: 0,
				TotalMoves:  0,
			}
			if err := gs.db.Create(&player).Error; err != nil {
				return nil, err
			}
		} else {
			return nil, result.Error
		}
	}

	// Deactivate existing active sessions
	gs.db.Model(&models.GameSession{}).
		Where("player_id = ? AND is_active = ?", player.ID, true).
		Update("is_active", false)

	// Initialize game data
	gameData := game.InitializeGameSession()

	// Create new game session
	gameSession := &models.GameSession{
		PlayerID:        player.ID,
		CurrentScore:    0,
		CurrentRow:      gameData["player_pos"].(map[string]int)["row"],
		CurrentCol:      gameData["player_pos"].(map[string]int)["col"],
		PreferredColumn: gameData["preferred_column"].(int),
		TotalMoves:      0,
		PearlsCollected: 0,
		IsActive:        true,
		IsCompleted:     false,
	}

	// Set game map and text grid
	gameSession.SetGameMap(gameData["game_map"].([][]int))
	gameSession.SetTextGrid(gameData["text_grid"].([][]string))

	if err := gs.db.Create(gameSession).Error; err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"success":       true,
		"session_token": gameSession.SessionToken,
		"game_data": map[string]interface{}{
			"text_grid":    gameData["text_grid"],
			"game_map":     gameSession.GetGameMap(),
			"player_pos":   map[string]int{"row": gameSession.CurrentRow, "col": gameSession.CurrentCol},
			"score":        gameSession.CurrentScore,
			"is_completed": gameSession.IsCompleted,
		},
	}, nil
}

// ProcessMove processes a move with full concurrency control
func (gs *GameService) ProcessMove(sessionToken, direction string) (map[string]interface{}, error) {
	// Get game session with lock
	var gameSession models.GameSession
	if err := gs.db.Where("session_token = ? AND is_active = ?", sessionToken, true).First(&gameSession).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return map[string]interface{}{
				"success": false,
				"error":   "Invalid or expired game session",
			}, nil
		}
		return nil, err
	}

	// Check if game is completed
	if gameSession.IsCompleted {
		return map[string]interface{}{
			"success": false,
			"error":   "Game already completed",
		}, nil
	}

	// Check game time limit
	if gs.isGameExpired(&gameSession) {
		gs.expireGame(&gameSession)
		return map[string]interface{}{
			"success": false,
			"error":   "Game expired due to time limit",
		}, nil
	}

	// Convert direction key to direction name
	directionName, exists := game.MovementKeys[direction]
	if !exists {
		return map[string]interface{}{
			"success": false,
			"error":   "Invalid movement key",
		}, nil
	}

	// Calculate new position
	gameMap := gameSession.GetGameMap()
	textGrid := gameSession.GetTextGrid()
	movementResult, err := game.CalculateNewPosition(
		directionName["direction"].(string),
		gameSession.CurrentRow,
		gameSession.CurrentCol,
		gameMap,
		textGrid,
		gameSession.PreferredColumn,
	)
	if err != nil {
		return map[string]interface{}{
			"success": false,
			"error":   err.Error(),
		}, nil
	}

	if !movementResult.IsValid {
		return map[string]interface{}{
			"success": false,
			"error":   "Out of bounds",
		}, nil
	}

	// Check if target position has a pearl
	pearlCollected := gameMap[movementResult.NewRow][movementResult.NewCol] == game.PEARL

	// Process the move atomically
	err = gs.db.Transaction(func(tx *gorm.DB) error {
		// Reload session in transaction to ensure fresh state
		var txGameSession models.GameSession
		if err := tx.Where("session_token = ?", sessionToken).First(&txGameSession).Error; err != nil {
			return err
		}

		// Process move with concurrency control
		err := txGameSession.ProcessMove(
			movementResult.NewRow,
			movementResult.NewCol,
			movementResult.PreferredColumn,
			pearlCollected,
			gs.cfg.PearlPoints,
		)
		if err != nil {
			return err
		}

		// Update game map
		updatedMap := txGameSession.GetGameMap()
		if pearlCollected {
			game.PlaceNewPearl(updatedMap, movementResult.NewRow, movementResult.NewCol)
			txGameSession.SetGameMap(updatedMap)
		}

		// Check if game should be completed
		if txGameSession.CurrentScore >= gs.cfg.TargetScore {
			txGameSession.CompleteGame()
			// Update player stats
			gs.updatePlayerStats(tx, txGameSession.PlayerID, &txGameSession)
		}

		// Validate score integrity
		if !txGameSession.ValidateScoreIntegrity(gs.cfg.PearlPoints) {
			return errors.New("score integrity validation failed")
		}

		// Save the session
		gameSession = txGameSession // Update our local copy
		return tx.Save(&txGameSession).Error
	})

	if err != nil {
		return map[string]interface{}{
			"success": false,
			"error":   err.Error(),
		}, nil
	}

	return map[string]interface{}{
		"success": true,
		"game_map": gameSession.GetGameMap(),
		"player_pos": map[string]int{
			"row": gameSession.CurrentRow,
			"col": gameSession.CurrentCol,
		},
		"score":           gameSession.CurrentScore,
		"pearl_collected": pearlCollected,
		"is_completed":    gameSession.IsCompleted,
		"completion_time": gameSession.CompletionTime,
		"final_score":     gameSession.FinalScore,
	}, nil
}

// GetGameState returns current game state
func (gs *GameService) GetGameState(sessionToken string) (map[string]interface{}, error) {
	var gameSession models.GameSession
	if err := gs.db.Where("session_token = ? AND is_active = ?", sessionToken, true).First(&gameSession).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return map[string]interface{}{
				"success": false,
				"error":   "Invalid or expired game session",
			}, nil
		}
		return nil, err
	}

	return map[string]interface{}{
		"success": true,
		"game_map": gameSession.GetGameMap(),
		"player_pos": map[string]int{
			"row": gameSession.CurrentRow,
			"col": gameSession.CurrentCol,
		},
		"score":            gameSession.CurrentScore,
		"is_completed":     gameSession.IsCompleted,
		"completion_time":  gameSession.CompletionTime,
		"final_score":      gameSession.FinalScore,
		"pearls_collected": gameSession.PearlsCollected,
		"total_moves":      gameSession.TotalMoves,
	}, nil
}

// GetLeaderboard returns leaderboard data
func (gs *GameService) GetLeaderboard(boardType string, limit int) (map[string]interface{}, error) {
	var sessions []models.GameSession
	query := gs.db.Preload("Player").Where("is_completed = ?", true)

	if boardType == "time" {
		query = query.Where("completion_time IS NOT NULL").Order("completion_time ASC")
	} else {
		query = query.Order("final_score DESC")
	}

	if err := query.Limit(limit).Find(&sessions).Error; err != nil {
		return nil, err
	}

	var leaderboard []map[string]interface{}
	for i, session := range sessions {
		entry := map[string]interface{}{
			"rank":             i + 1,
			"username":         session.Player.Username,
			"score":            session.FinalScore,
			"completion_time":  session.CompletionTime,
			"total_moves":      session.TotalMoves,
			"pearls_collected": session.PearlsCollected,
		}
		if session.EndTime != nil {
			entry["completed_at"] = session.EndTime.Format(time.RFC3339)
		}
		leaderboard = append(leaderboard, entry)
	}

	return map[string]interface{}{
		"success":     true,
		"leaderboard": leaderboard,
		"type":        boardType,
	}, nil
}

// Helper methods
func (gs *GameService) isGameExpired(gameSession *models.GameSession) bool {
	if gameSession.StartTime == nil {
		return false
	}
	return time.Since(*gameSession.StartTime) > gs.cfg.MaxGameTime
}

func (gs *GameService) expireGame(gameSession *models.GameSession) {
	gameSession.IsActive = false
	now := time.Now()
	gameSession.EndTime = &now
	gs.db.Save(gameSession)
}

func (gs *GameService) updatePlayerStats(tx *gorm.DB, playerID uint, gameSession *models.GameSession) {
	var player models.Player
	tx.First(&player, playerID)

	player.TotalGames++
	if gameSession.IsCompleted {
		player.CompletedGames++
		if gameSession.FinalScore != nil && *gameSession.FinalScore > player.BestScore {
			player.BestScore = *gameSession.FinalScore
		}
		if gameSession.CompletionTime != nil && (player.FastestTime == nil || *gameSession.CompletionTime < *player.FastestTime) {
			player.FastestTime = gameSession.CompletionTime
		}
	}
	player.TotalPearls += gameSession.PearlsCollected
	player.TotalMoves += gameSession.TotalMoves

	tx.Save(&player)
}