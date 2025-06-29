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
	"w": {"direction": "word_forward", "description": "Move forward to beginning of next word"},
	"W": {"direction": "word_forward_space", "description": "Move forward to beginning of next WORD (space-separated)"},
	"b": {"direction": "word_backward", "description": "Move backward to beginning of current/previous word"},
	"B": {"direction": "word_backward_space", "description": "Move backward to beginning of current/previous WORD"},
	"e": {"direction": "word_end", "description": "Move to end of current/next word"},
	"$": {"direction": "line_end", "description": "Move to end of current line"},
	"0": {"direction": "line_start", "description": "Move to beginning of current line"},
}

// ValidMovementKeys list of all valid movement keys
var ValidMovementKeys = []string{"h", "j", "k", "l", "w", "W", "b", "B", "e", "$", "0"}

// MovementResult represents the result of a movement calculation
type MovementResult struct {
	NewRow          int  `json:"new_row"`
	NewCol          int  `json:"new_col"`
	PreferredColumn int  `json:"preferred_column"`
	IsValid         bool `json:"is_valid"`
}

// CalculateNewPosition calculates the new position based on vim-style movement
func CalculateNewPosition(direction string, currentRow, currentCol int, gameMap [][]int, textGrid [][]string, preferredColumn int) (*MovementResult, error) {
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
		newCol = clampToRow(preferredColumn, newRow, gameMap)
	case "down":
		newRow = currentRow + 1
		newCol = clampToRow(preferredColumn, newRow, gameMap)
	case "word_forward", "word_forward_space":
		newRow, newCol = findWordForward(currentRow, currentCol, textGrid, direction == "word_forward_space")
		newPreferredColumn = newCol
	case "word_backward", "word_backward_space":
		newRow, newCol = findWordBackward(currentRow, currentCol, textGrid, direction == "word_backward_space")
		newPreferredColumn = newCol
	case "word_end":
		newRow, newCol = findWordEnd(currentRow, currentCol, textGrid)
		newPreferredColumn = newCol
	case "line_end":
		newCol = len(gameMap[currentRow]) - 1
		newPreferredColumn = newCol
	case "line_start":
		newCol = 0
		newPreferredColumn = newCol
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
		"left":               true,
		"right":              true,
		"up":                 true,
		"down":               true,
		"word_forward":       true,
		"word_forward_space": true,
		"word_backward":      true,
		"word_backward_space": true,
		"word_end":           true,
		"line_end":           true,
		"line_start":         true,
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

// Helper functions for vim-like word movement

// isWordChar determines if a character is a word character
func isWordChar(char string) bool {
	if len(char) == 0 {
		return false
	}
	c := char[0]
	// ASCII letters, digits, underscore
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') || c == '_'
}

// isSpace determines if a character is whitespace
func isSpace(char string) bool {
	return char == " " || char == "\t"
}

// isPunctuation determines if a character is punctuation
func isPunctuation(char string) bool {
	if len(char) == 0 {
		return false
	}
	return !isWordChar(char) && !isSpace(char)
}

func findWordForward(row, col int, textGrid [][]string, spaceSeparated bool) (int, int) {
	if row < 0 || row >= len(textGrid) {
		return row, col
	}
	
	currentRow := row
	currentCol := col
	
	// Skip current word/punctuation
	for currentRow < len(textGrid) && currentCol < len(textGrid[currentRow]) {
		currentChar := textGrid[currentRow][currentCol]
		if isSpace(currentChar) {
			break
		}
		currentCol++
	}
	
	// Skip whitespace
	for currentRow < len(textGrid) && currentCol < len(textGrid[currentRow]) {
		currentChar := textGrid[currentRow][currentCol]
		if !isSpace(currentChar) {
			break
		}
		currentCol++
	}
	
	// If we reached end of line, go to next line
	if currentRow < len(textGrid) && currentCol >= len(textGrid[currentRow]) {
		currentRow++
		currentCol = 0
		// Skip leading whitespace on new line
		for currentRow < len(textGrid) && currentCol < len(textGrid[currentRow]) {
			currentChar := textGrid[currentRow][currentCol]
			if !isSpace(currentChar) {
				break
			}
			currentCol++
		}
	}
	
	// Ensure bounds
	if currentRow >= len(textGrid) {
		currentRow = len(textGrid) - 1
		currentCol = len(textGrid[currentRow]) - 1
	} else if currentRow >= 0 && currentCol >= len(textGrid[currentRow]) {
		currentCol = len(textGrid[currentRow]) - 1
	}
	
	return currentRow, currentCol
}

func findWordBackward(row, col int, textGrid [][]string, spaceSeparated bool) (int, int) {
	if row < 0 || row >= len(textGrid) {
		return row, col
	}
	
	currentRow := row
	currentCol := col
	
	// Move one position back to start
	currentCol--
	if currentCol < 0 {
		if currentRow > 0 {
			currentRow--
			if currentRow >= 0 && currentRow < len(textGrid) {
				currentCol = len(textGrid[currentRow]) - 1
			}
		} else {
			return 0, 0 // At beginning
		}
	}
	
	if currentRow < 0 || currentCol < 0 {
		return 0, 0
	}
	
	// Skip whitespace backwards
	for currentRow >= 0 && currentCol >= 0 {
		if currentRow >= len(textGrid) || currentCol >= len(textGrid[currentRow]) {
			break
		}
		currentChar := textGrid[currentRow][currentCol]
		if !isSpace(currentChar) {
			break
		}
		currentCol--
		if currentCol < 0 && currentRow > 0 {
			currentRow--
			if currentRow >= 0 && currentRow < len(textGrid) {
				currentCol = len(textGrid[currentRow]) - 1
			}
		}
	}
	
	// Find beginning of current word
	if currentRow >= 0 && currentCol >= 0 && currentRow < len(textGrid) && currentCol < len(textGrid[currentRow]) {
		currentChar := textGrid[currentRow][currentCol]
		isCurrentWord := isWordChar(currentChar)
		isCurrentPunct := isPunctuation(currentChar)
		
		// Move to beginning of current word/punctuation group
		for currentCol > 0 {
			prevChar := textGrid[currentRow][currentCol-1]
			
			if spaceSeparated {
				// For B: only stop at whitespace
				if isSpace(prevChar) {
					break
				}
			} else {
				// For b: stop when character type changes
				if isSpace(prevChar) || 
				   (isCurrentWord && !isWordChar(prevChar)) ||
				   (isCurrentPunct && !isPunctuation(prevChar)) {
					break
				}
			}
			currentCol--
		}
	}
	
	// Ensure bounds
	if currentRow < 0 {
		currentRow = 0
		currentCol = 0
	} else if currentCol < 0 {
		currentCol = 0
	}
	
	return currentRow, currentCol
}

func findWordEnd(row, col int, textGrid [][]string) (int, int) {
	if row < 0 || row >= len(textGrid) {
		return row, col
	}
	
	currentRow := row
	currentCol := col
	
	// If we're at end of word, move to next word first
	if currentCol < len(textGrid[currentRow]) {
		currentChar := textGrid[currentRow][currentCol]
		if !isSpace(currentChar) {
			// Skip current word/punctuation to find next
			isCurrentWord := isWordChar(currentChar)
			isCurrentPunct := isPunctuation(currentChar)
			
			for currentRow < len(textGrid) && currentCol < len(textGrid[currentRow]) {
				curChar := textGrid[currentRow][currentCol]
				if isSpace(curChar) ||
				   (isCurrentWord && !isWordChar(curChar)) ||
				   (isCurrentPunct && !isPunctuation(curChar)) {
					break
				}
				currentCol++
			}
			
			// Skip whitespace to next word
			for currentRow < len(textGrid) && currentCol < len(textGrid[currentRow]) {
				curChar := textGrid[currentRow][currentCol]
				if !isSpace(curChar) {
					break
				}
				currentCol++
			}
			
			// Handle end of line
			if currentRow < len(textGrid) && currentCol >= len(textGrid[currentRow]) {
				currentRow++
				currentCol = 0
				// Skip leading whitespace
				for currentRow < len(textGrid) && currentCol < len(textGrid[currentRow]) {
					curChar := textGrid[currentRow][currentCol]
					if !isSpace(curChar) {
						break
					}
					currentCol++
				}
			}
		}
	}
	
	// Now find end of current word
	if currentRow < len(textGrid) && currentCol < len(textGrid[currentRow]) {
		currentChar := textGrid[currentRow][currentCol]
		isCurrentWord := isWordChar(currentChar)
		isCurrentPunct := isPunctuation(currentChar)
		
		// Move to end of current word/punctuation
		for currentCol < len(textGrid[currentRow])-1 {
			nextChar := textGrid[currentRow][currentCol+1]
			if isSpace(nextChar) ||
			   (isCurrentWord && !isWordChar(nextChar)) ||
			   (isCurrentPunct && !isPunctuation(nextChar)) {
				break
			}
			currentCol++
		}
	}
	
	// Ensure bounds
	if currentRow >= len(textGrid) {
		currentRow = len(textGrid) - 1
		currentCol = len(textGrid[currentRow]) - 1
	} else if currentRow >= 0 && currentCol >= len(textGrid[currentRow]) {
		currentCol = len(textGrid[currentRow]) - 1
	}
	
	return currentRow, currentCol
}
