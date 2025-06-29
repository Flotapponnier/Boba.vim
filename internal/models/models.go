package models

import (
	"encoding/json"
	"sync"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Player struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Username  string    `gorm:"unique;not null" json:"username"`
	Email     string    `gorm:"unique;not null" json:"email"`
	Password  string    `gorm:"not null" json:"-"` // Don't include in JSON responses
	IsRegistered bool   `gorm:"default:false" json:"is_registered"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Stats
	TotalGames       int `json:"total_games"`
	CompletedGames   int `json:"completed_games"`
	BestScore        int `json:"best_score"`
	TotalPearls      int `json:"total_pearls"`
	TotalMoves       int `json:"total_moves"`
	FastestTime      *int `json:"fastest_time"`
}

type GameSession struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	SessionToken  string    `gorm:"unique;not null" json:"session_token"`
	PlayerID      uint      `json:"player_id"`
	Player        Player    `gorm:"foreignKey:PlayerID" json:"player,omitempty"`
	SelectedCharacter string `gorm:"default:boba" json:"selected_character"`
	
	// Game state with mutex for concurrent access
	gameMapMutex  sync.RWMutex `gorm:"-" json:"-"`
	GameMapJSON   string       `json:"-"`
	gameMap       [][]int      `gorm:"-" json:"game_map"`
	
	// Text grid for movement calculations
	textGridMutex sync.RWMutex  `gorm:"-" json:"-"`
	TextGridJSON  string        `json:"-"`
	textGrid      [][]string    `gorm:"-" json:"text_grid"`
	
	// Position and game state
	CurrentRow      int  `json:"current_row"`
	CurrentCol      int  `json:"current_col"`
	PreferredColumn int  `json:"preferred_column"`
	CurrentScore    int  `json:"current_score"`
	FinalScore      *int `json:"final_score"`
	
	// Move tracking with mutex
	moveMutex     sync.Mutex `gorm:"-" json:"-"`
	TotalMoves    int        `json:"total_moves"`
	PearlsCollected int      `json:"pearls_collected"`
	LastMoveTime  *time.Time `json:"last_move_time"`
	
	// Game status
	IsActive    bool       `json:"is_active"`
	IsCompleted bool       `json:"is_completed"`
	StartTime   *time.Time `json:"start_time"`
	EndTime     *time.Time `json:"end_time"`
	CompletionTime *int    `json:"completion_time"` // seconds
	
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// BeforeCreate sets session token and start time
func (gs *GameSession) BeforeCreate(tx *gorm.DB) error {
	gs.SessionToken = uuid.New().String()
	now := time.Now()
	gs.StartTime = &now
	return nil
}

// AfterFind loads game map and text grid from JSON
func (gs *GameSession) AfterFind(tx *gorm.DB) error {
	if gs.GameMapJSON != "" {
		if err := json.Unmarshal([]byte(gs.GameMapJSON), &gs.gameMap); err != nil {
			return err
		}
	}
	if gs.TextGridJSON != "" {
		if err := json.Unmarshal([]byte(gs.TextGridJSON), &gs.textGrid); err != nil {
			return err
		}
	}
	return nil
}

// BeforeSave saves game map and text grid to JSON
func (gs *GameSession) BeforeSave(tx *gorm.DB) error {
	gs.gameMapMutex.RLock()
	defer gs.gameMapMutex.RUnlock()
	
	if gs.gameMap != nil {
		mapJSON, err := json.Marshal(gs.gameMap)
		if err != nil {
			return err
		}
		gs.GameMapJSON = string(mapJSON)
	}
	
	gs.textGridMutex.RLock()
	defer gs.textGridMutex.RUnlock()
	
	if gs.textGrid != nil {
		textJSON, err := json.Marshal(gs.textGrid)
		if err != nil {
			return err
		}
		gs.TextGridJSON = string(textJSON)
	}
	return nil
}

// GetGameMap returns a copy of the game map safely
func (gs *GameSession) GetGameMap() [][]int {
	gs.gameMapMutex.RLock()
	defer gs.gameMapMutex.RUnlock()
	
	if gs.gameMap == nil {
		return nil
	}
	
	// Return deep copy to prevent external modification
	mapCopy := make([][]int, len(gs.gameMap))
	for i, row := range gs.gameMap {
		mapCopy[i] = make([]int, len(row))
		copy(mapCopy[i], row)
	}
	return mapCopy
}

// SetGameMap sets the game map safely
func (gs *GameSession) SetGameMap(gameMap [][]int) {
	gs.gameMapMutex.Lock()
	defer gs.gameMapMutex.Unlock()
	
	// Create deep copy to prevent external modification
	gs.gameMap = make([][]int, len(gameMap))
	for i, row := range gameMap {
		gs.gameMap[i] = make([]int, len(row))
		copy(gs.gameMap[i], row)
	}
}

// GetTextGrid returns a copy of the text grid safely
func (gs *GameSession) GetTextGrid() [][]string {
	gs.textGridMutex.RLock()
	defer gs.textGridMutex.RUnlock()
	
	if gs.textGrid == nil {
		return nil
	}
	
	// Return deep copy to prevent external modification
	textCopy := make([][]string, len(gs.textGrid))
	for i, row := range gs.textGrid {
		textCopy[i] = make([]string, len(row))
		copy(textCopy[i], row)
	}
	return textCopy
}

// SetTextGrid sets the text grid safely
func (gs *GameSession) SetTextGrid(textGrid [][]string) {
	gs.textGridMutex.Lock()
	defer gs.textGridMutex.Unlock()
	
	// Create deep copy to prevent external modification
	gs.textGrid = make([][]string, len(textGrid))
	for i, row := range textGrid {
		gs.textGrid[i] = make([]string, len(row))
		copy(gs.textGrid[i], row)
	}
}

// ProcessMove handles a move with proper concurrency control
func (gs *GameSession) ProcessMove(newRow, newCol, preferredCol int, pearlCollected bool, pearlPoints int) error {
	gs.moveMutex.Lock()
	defer gs.moveMutex.Unlock()
	
	now := time.Now()
	
	// Check if enough time has passed since last move (prevent spam)
	if gs.LastMoveTime != nil && now.Sub(*gs.LastMoveTime) < 50*time.Millisecond {
		return ErrMoveTooFast
	}
	
	gs.gameMapMutex.Lock()
	defer gs.gameMapMutex.Unlock()
	
	// Update map
	if gs.gameMap != nil {
		gs.gameMap[gs.CurrentRow][gs.CurrentCol] = 0
		gs.gameMap[newRow][newCol] = 1
	}
	
	// Update position
	gs.CurrentRow = newRow
	gs.CurrentCol = newCol
	gs.PreferredColumn = preferredCol
	gs.TotalMoves++
	gs.LastMoveTime = &now
	
	// Handle pearl collection
	if pearlCollected {
		gs.CurrentScore += pearlPoints
		gs.PearlsCollected++
	}
	
	return nil
}

// CompleteGame marks the game as completed
func (gs *GameSession) CompleteGame() {
	gs.moveMutex.Lock()
	defer gs.moveMutex.Unlock()
	
	now := time.Now()
	gs.IsCompleted = true
	gs.IsActive = false
	gs.EndTime = &now
	gs.FinalScore = &gs.CurrentScore
	
	if gs.StartTime != nil {
		completionTime := int(now.Sub(*gs.StartTime).Seconds())
		gs.CompletionTime = &completionTime
	}
}

// ValidateScoreIntegrity validates that the score matches pearl collection
func (gs *GameSession) ValidateScoreIntegrity(pearlPoints int) bool {
	expectedScore := gs.PearlsCollected * pearlPoints
	return gs.CurrentScore == expectedScore
}

// Custom errors
var (
	ErrMoveTooFast = &GameError{Code: "MOVE_TOO_FAST", Message: "Move requests too frequent"}
	ErrGameCompleted = &GameError{Code: "GAME_COMPLETED", Message: "Game already completed"}
	ErrInvalidMove = &GameError{Code: "INVALID_MOVE", Message: "Invalid move"}
)

type GameError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (e *GameError) Error() string {
	return e.Message
}