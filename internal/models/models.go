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
	
	// Game state with mutex for concurrent access
	gameMapMutex  sync.RWMutex `gorm:"-" json:"-"`
	GameMapJSON   string       `json:"-"`
	gameMap       [][]int      `gorm:"-" json:"game_map"`
	
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

// AfterFind loads game map from JSON
func (gs *GameSession) AfterFind(tx *gorm.DB) error {
	if gs.GameMapJSON != "" {
		return json.Unmarshal([]byte(gs.GameMapJSON), &gs.gameMap)
	}
	return nil
}

// BeforeSave saves game map to JSON
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