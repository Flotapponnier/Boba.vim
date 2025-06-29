package game

import (
	"errors"
	"fmt"
)

// Movement directions
var MovementKeys = map[string]map[string]interface{}{
	"h": {"direction": "left", "description": "Move left (vim style)"},
	"j": {"direction": "down", "description": "Move down (vim style)"},
	"k": {"direction": "up", "description": "Move up (vim style)"},
	"l": {"direction": "right", "description": "Move right (vim style)"},
}

// ValidMovementKeys list of all valid movement keys
var ValidMovementKeys = []string{"h", "j", "k", "l"}

// MovementResult represents the result of a movement calculation
type MovementResult struct {
	NewRow          int  `json:"new_row"`
	NewCol          int  `json:"new_col"`
	PreferredColumn int  `json:"preferred_column"`
	IsValid         bool `json:"is_valid"`
}

// CalculateNewPosition calculates the new position based on vim-style movement
func CalculateNewPosition(direction string, currentRow, currentCol int, gameMap [][]int, preferredColumn int) (*MovementResult, error) {
	if !isValidDirection(direction) {
		return nil, fmt.Errorf("invalid direction: %s", direction)
	}
	
	newRow, newCol := currentRow, currentCol
	newPreferredColumn := preferredColumn
	
	switch direction {
	case "left":
		newCol = currentCol - 1
		newPreferredColumn = newCol
	case "right":
		newCol = currentCol + 1
		newPreferredColumn = newCol
	case "up":
		newRow = currentRow - 1
		// Try to maintain preferred column, but clamp to valid range
		newCol = clampToRow(preferredColumn, newRow, gameMap)
	case "down":
		newRow = currentRow + 1
		// Try to maintain preferred column, but clamp to valid range
		newCol = clampToRow(preferredColumn, newRow, gameMap)
	default:
		return nil, fmt.Errorf("unknown direction: %s", direction)
	}
	
	isValid := IsValidPosition(newRow, newCol, gameMap)
	
	return &MovementResult{
		NewRow:          newRow,
		NewCol:          newCol,
		PreferredColumn: newPreferredColumn,
		IsValid:         isValid,
	}, nil
}

// isValidDirection checks if the direction is valid
func isValidDirection(direction string) bool {
	validDirections := map[string]bool{
		"left":  true,
		"right": true,
		"up":    true,
		"down":  true,
	}
	return validDirections[direction]
}

// clampToRow clamps the column to valid range for the given row
func clampToRow(preferredCol, row int, gameMap [][]int) int {
	if row < 0 || row >= len(gameMap) {
		return preferredCol
	}
	
	maxCol := len(gameMap[row]) - 1
	if preferredCol < 0 {
		return 0
	}
	if preferredCol > maxCol {
		return maxCol
	}
	return preferredCol
}

// GetAvailableMovements returns all available movement keys
func GetAvailableMovements() []map[string]interface{} {
	var movements []map[string]interface{}
	
	for key, info := range MovementKeys {
		movement := map[string]interface{}{
			"key":         key,
			"direction":   info["direction"],
			"description": info["description"],
		}
		movements = append(movements, movement)
	}
	
	return movements
}

// ValidateMovement validates if a movement is allowed
func ValidateMovement(direction string, currentRow, currentCol int, gameMap [][]int) error {
	if !isValidDirection(direction) {
		return errors.New("invalid direction")
	}
	
	if currentRow < 0 || currentRow >= len(gameMap) || 
	   currentCol < 0 || currentCol >= len(gameMap[0]) {
		return errors.New("current position out of bounds")
	}
	
	return nil
}