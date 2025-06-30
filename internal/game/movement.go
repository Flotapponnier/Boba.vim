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
	"E": {"direction": "word_end_space", "description": "Move to end of current/next WORD (space-separated)"},
	"$": {"direction": "line_end", "description": "Move to end of current line"},
	"0": {"direction": "line_start", "description": "Move to beginning of current line"},
	"^": {"direction": "line_first_non_blank", "description": "Move to first non-blank character of line"},
	"g_": {"direction": "line_last_non_blank", "description": "Move to last non-blank character of line"},
	"gg": {"direction": "file_start", "description": "Go to top of file"},
	"G": {"direction": "file_end", "description": "Go to bottom of file"},
	"H": {"direction": "screen_top", "description": "Go to top of screen"},
	"M": {"direction": "screen_middle", "description": "Go to middle of screen"},
	"L": {"direction": "screen_bottom", "description": "Go to bottom of screen"},
	"{": {"direction": "paragraph_prev", "description": "Go to previous paragraph"},
	"}": {"direction": "paragraph_next", "description": "Go to next paragraph"},
	"(": {"direction": "sentence_prev", "description": "Go to previous sentence"},
	")": {"direction": "sentence_next", "description": "Go to next sentence"},
	"f": {"direction": "find_char_forward", "description": "Find character forward"},
	"F": {"direction": "find_char_backward", "description": "Find character backward"},
	"t": {"direction": "till_char_forward", "description": "Till character forward"},
	"T": {"direction": "till_char_backward", "description": "Till character backward"},
}

// ValidMovementKeys list of all valid movement keys
var ValidMovementKeys = []string{"h", "j", "k", "l", "w", "W", "b", "B", "e", "E", "$", "0", "^", "g_", "gg", "G", "H", "M", "L", "{", "}", "(", ")", "f", "F", "t", "T"}

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
	case "word_end_space":
		newRow, newCol = findWordEndSpace(currentRow, currentCol, textGrid)
		newPreferredColumn = newCol
	case "line_end":
		newCol = len(gameMap[currentRow]) - 1
		newPreferredColumn = newCol
	case "line_start":
		newCol = 0
		newPreferredColumn = newCol
	case "line_first_non_blank":
		newCol = findFirstNonBlank(currentRow, textGrid)
		newPreferredColumn = newCol
	case "line_last_non_blank":
		newCol = findLastNonBlank(currentRow, textGrid)
		newPreferredColumn = newCol
	case "file_start":
		newRow = 0
		newCol = 0
		newPreferredColumn = newCol
	case "file_end":
		newRow = len(gameMap) - 1
		newCol = 0
		newPreferredColumn = newCol
	case "screen_top":
		newRow = 0
		newCol = clampToRow(preferredColumn, newRow, gameMap)
	case "screen_middle":
		newRow = len(gameMap) / 2
		newCol = clampToRow(preferredColumn, newRow, gameMap)
	case "screen_bottom":
		newRow = len(gameMap) - 1
		newCol = clampToRow(preferredColumn, newRow, gameMap)
	case "paragraph_prev":
		newRow, newCol = findParagraphPrev(currentRow, currentCol, textGrid)
		newPreferredColumn = newCol
	case "paragraph_next":
		newRow, newCol = findParagraphNext(currentRow, currentCol, textGrid)
		newPreferredColumn = newCol
	case "sentence_prev":
		newRow, newCol = findSentencePrev(currentRow, currentCol, textGrid)
		newPreferredColumn = newCol
	case "sentence_next":
		newRow, newCol = findSentenceNext(currentRow, currentCol, textGrid)
		newPreferredColumn = newCol
	default:
		// Check if it's a character search direction with character parameter
		if len(direction) > 17 && direction[:17] == "find_char_forward" {
			// Extract character from direction string (format: "find_char_forward_X")
			if len(direction) == 19 && direction[17] == '_' {
				targetChar := string(direction[18])
				newRow, newCol = findCharForward(currentRow, currentCol, textGrid, targetChar)
				newPreferredColumn = newCol
			} else {
				return nil, fmt.Errorf("invalid find_char_forward format: %s", direction)
			}
		} else if len(direction) > 18 && direction[:18] == "find_char_backward" {
			// Extract character from direction string (format: "find_char_backward_X")
			if len(direction) == 20 && direction[18] == '_' {
				targetChar := string(direction[19])
				newRow, newCol = findCharBackward(currentRow, currentCol, textGrid, targetChar)
				newPreferredColumn = newCol
			} else {
				return nil, fmt.Errorf("invalid find_char_backward format: %s", direction)
			}
		} else if len(direction) > 17 && direction[:17] == "till_char_forward" {
			// Extract character from direction string (format: "till_char_forward_X")
			if len(direction) == 19 && direction[17] == '_' {
				targetChar := string(direction[18])
				newRow, newCol = tillCharForward(currentRow, currentCol, textGrid, targetChar)
				newPreferredColumn = newCol
			} else {
				return nil, fmt.Errorf("invalid till_char_forward format: %s", direction)
			}
		} else if len(direction) > 18 && direction[:18] == "till_char_backward" {
			// Extract character from direction string (format: "till_char_backward_X")
			if len(direction) == 20 && direction[18] == '_' {
				targetChar := string(direction[19])
				newRow, newCol = tillCharBackward(currentRow, currentCol, textGrid, targetChar)
				newPreferredColumn = newCol
			} else {
				return nil, fmt.Errorf("invalid till_char_backward format: %s", direction)
			}
		} else {
			return nil, fmt.Errorf("unknown direction: %s", direction)
		}
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
		"left":                   true,
		"right":                  true,
		"up":                     true,
		"down":                   true,
		"word_forward":           true,
		"word_forward_space":     true,
		"word_backward":          true,
		"word_backward_space":    true,
		"word_end":               true,
		"word_end_space":         true,
		"line_end":               true,
		"line_start":             true,
		"line_first_non_blank":   true,
		"line_last_non_blank":    true,
		"file_start":             true,
		"file_end":               true,
		"screen_top":             true,
		"screen_middle":          true,
		"screen_bottom":          true,
		"paragraph_prev":         true,
		"paragraph_next":         true,
		"sentence_prev":          true,
		"sentence_next":          true,
		"find_char_forward":      true,
		"find_char_backward":     true,
		"till_char_forward":      true,
		"till_char_backward":     true,
	}
	
	// Check standard directions first
	if validDirections[direction] {
		return true
	}
	
	// Check character search directions with character parameter
	if len(direction) == 19 && direction[:17] == "find_char_forward" && direction[17] == '_' {
		return true
	}
	if len(direction) == 20 && direction[:18] == "find_char_backward" && direction[18] == '_' {
		return true
	}
	if len(direction) == 19 && direction[:17] == "till_char_forward" && direction[17] == '_' {
		return true
	}
	if len(direction) == 20 && direction[:18] == "till_char_backward" && direction[18] == '_' {
		return true
	}
	
	return false
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


func findWordEndSpace(row, col int, textGrid [][]string) (int, int) {
	if row < 0 || row >= len(textGrid) {
		return row, col
	}
	
	currentRow := row
	currentCol := col
	
	// If we're at end of WORD, move to next WORD first
	if currentCol < len(textGrid[currentRow]) {
		currentChar := textGrid[currentRow][currentCol]
		if !isSpace(currentChar) {
			// Skip current WORD (everything until space)
			for currentRow < len(textGrid) && currentCol < len(textGrid[currentRow]) {
				curChar := textGrid[currentRow][currentCol]
				if isSpace(curChar) {
					break
				}
				currentCol++
			}
			
			// Skip whitespace to next WORD
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
	
	// Now find end of current WORD (space-separated)
	if currentRow < len(textGrid) && currentCol < len(textGrid[currentRow]) {
		// Move to end of current WORD
		for currentCol < len(textGrid[currentRow])-1 {
			nextChar := textGrid[currentRow][currentCol+1]
			if isSpace(nextChar) {
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

func findFirstNonBlank(row int, textGrid [][]string) int {
	if row < 0 || row >= len(textGrid) {
		return 0
	}
	
	// Find first non-blank character in the row
	for col := 0; col < len(textGrid[row]); col++ {
		if !isSpace(textGrid[row][col]) {
			return col
		}
	}
	
	// If entire row is blank, return 0
	return 0
}

func findLastNonBlank(row int, textGrid [][]string) int {
	if row < 0 || row >= len(textGrid) {
		return 0
	}
	
	// Find last non-blank character in the row
	for col := len(textGrid[row]) - 1; col >= 0; col-- {
		if !isSpace(textGrid[row][col]) {
			return col
		}
	}
	
	// If entire row is blank, return 0
	return 0
}

func findParagraphPrev(row, col int, textGrid [][]string) (int, int) {
	if row < 0 || row >= len(textGrid) {
		return row, col
	}
	
	currentRow := row
	
	// Move up from current position
	currentRow--
	
	// Skip current paragraph (non-empty lines)
	for currentRow >= 0 {
		// Check if line is empty (only whitespace)
		isEmpty := true
		for _, char := range textGrid[currentRow] {
			if !isSpace(char) {
				isEmpty = false
				break
			}
		}
		if isEmpty {
			break
		}
		currentRow--
	}
	
	// Skip empty lines
	for currentRow >= 0 {
		isEmpty := true
		for _, char := range textGrid[currentRow] {
			if !isSpace(char) {
				isEmpty = false
				break
			}
		}
		if !isEmpty {
			break
		}
		currentRow--
	}
	
	// If we found a paragraph, go to its beginning
	if currentRow >= 0 {
		// Move to beginning of this paragraph
		for currentRow > 0 {
			prevRowEmpty := true
			for _, char := range textGrid[currentRow-1] {
				if !isSpace(char) {
					prevRowEmpty = false
					break
				}
			}
			if prevRowEmpty {
				break
			}
			currentRow--
		}
	} else {
		currentRow = 0
	}
	
	return currentRow, 0
}

func findParagraphNext(row, col int, textGrid [][]string) (int, int) {
	if row < 0 || row >= len(textGrid) {
		return row, col
	}
	
	currentRow := row
	
	// Move down from current position
	currentRow++
	
	// Skip current paragraph (non-empty lines)
	for currentRow < len(textGrid) {
		// Check if line is empty (only whitespace)
		isEmpty := true
		for _, char := range textGrid[currentRow] {
			if !isSpace(char) {
				isEmpty = false
				break
			}
		}
		if isEmpty {
			break
		}
		currentRow++
	}
	
	// Skip empty lines
	for currentRow < len(textGrid) {
		isEmpty := true
		for _, char := range textGrid[currentRow] {
			if !isSpace(char) {
				isEmpty = false
				break
			}
		}
		if !isEmpty {
			break
		}
		currentRow++
	}
	
	// Ensure bounds
	if currentRow >= len(textGrid) {
		currentRow = len(textGrid) - 1
	}
	
	return currentRow, 0
}

func findSentencePrev(row, col int, textGrid [][]string) (int, int) {
	if row < 0 || row >= len(textGrid) {
		return row, col
	}
	
	currentRow := row
	currentCol := col
	
	// Move backward to find sentence start
	for {
		// Move one position back
		currentCol--
		if currentCol < 0 {
			if currentRow > 0 {
				currentRow--
				if currentRow >= 0 && currentRow < len(textGrid) {
					currentCol = len(textGrid[currentRow]) - 1
				}
			} else {
				return 0, 0 // At beginning of file
			}
		}
		
		if currentRow < 0 || currentCol < 0 {
			return 0, 0
		}
		
		// Check if current character ends a sentence
		if currentRow < len(textGrid) && currentCol < len(textGrid[currentRow]) {
			char := textGrid[currentRow][currentCol]
			if char == "." || char == "!" || char == "?" {
				// Found sentence end, move to start of next sentence
				currentCol++
				if currentCol >= len(textGrid[currentRow]) {
					if currentRow+1 < len(textGrid) {
						currentRow++
						currentCol = 0
					}
				}
				
				// Skip whitespace to sentence start
				for currentRow < len(textGrid) && currentCol < len(textGrid[currentRow]) {
					if !isSpace(textGrid[currentRow][currentCol]) {
						break
					}
					currentCol++
				}
				
				return currentRow, currentCol
			}
		}
		
		// Prevent infinite loop
		if currentRow == 0 && currentCol == 0 {
			break
		}
	}
	
	return 0, 0
}

func findSentenceNext(row, col int, textGrid [][]string) (int, int) {
	if row < 0 || row >= len(textGrid) {
		return row, col
	}
	
	currentRow := row
	currentCol := col
	
	// Move forward to find sentence end
	for currentRow < len(textGrid) {
		for currentCol < len(textGrid[currentRow]) {
			char := textGrid[currentRow][currentCol]
			if char == "." || char == "!" || char == "?" {
				// Found sentence end, move to start of next sentence
				currentCol++
				if currentCol >= len(textGrid[currentRow]) {
					if currentRow+1 < len(textGrid) {
						currentRow++
						currentCol = 0
					}
				}
				
				// Skip whitespace to next sentence start
				for currentRow < len(textGrid) && currentCol < len(textGrid[currentRow]) {
					if !isSpace(textGrid[currentRow][currentCol]) {
						break
					}
					currentCol++
				}
				
				// Handle end of line
				if currentRow < len(textGrid) && currentCol >= len(textGrid[currentRow]) {
					if currentRow+1 < len(textGrid) {
						currentRow++
						currentCol = 0
						// Skip leading whitespace
						for currentRow < len(textGrid) && currentCol < len(textGrid[currentRow]) {
							if !isSpace(textGrid[currentRow][currentCol]) {
								break
							}
							currentCol++
						}
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
			currentCol++
		}
		currentRow++
		currentCol = 0
	}
	
	// If no sentence end found, go to end of file
	if len(textGrid) > 0 {
		return len(textGrid) - 1, len(textGrid[len(textGrid)-1]) - 1
	}
	return row, col
}

func findCharForward(row, col int, textGrid [][]string, targetChar string) (int, int) {
	if row < 0 || row >= len(textGrid) {
		return row, col
	}
	
	// Search only in current row, starting from next position (vim behavior)
	for searchCol := col + 1; searchCol < len(textGrid[row]); searchCol++ {
		if textGrid[row][searchCol] == targetChar {
			return row, searchCol
		}
	}
	
	// Character not found in current line, return original position
	return row, col
}

func findCharBackward(row, col int, textGrid [][]string, targetChar string) (int, int) {
	if row < 0 || row >= len(textGrid) {
		return row, col
	}
	
	// Search only in current row (backwards), starting from previous position (vim behavior)
	for searchCol := col - 1; searchCol >= 0; searchCol-- {
		if textGrid[row][searchCol] == targetChar {
			return row, searchCol
		}
	}
	
	// Character not found in current line, return original position
	return row, col
}

func tillCharForward(row, col int, textGrid [][]string, targetChar string) (int, int) {
	if row < 0 || row >= len(textGrid) {
		return row, col
	}
	
	// Search only in current row, starting from next position (vim behavior)
	for searchCol := col + 1; searchCol < len(textGrid[row]); searchCol++ {
		if textGrid[row][searchCol] == targetChar {
			// Return position one before the target
			return row, searchCol - 1
		}
	}
	
	// Character not found in current line, return original position
	return row, col
}

func tillCharBackward(row, col int, textGrid [][]string, targetChar string) (int, int) {
	if row < 0 || row >= len(textGrid) {
		return row, col
	}
	
	// Search only in current row (backwards), starting from previous position (vim behavior)
	for searchCol := col - 1; searchCol >= 0; searchCol-- {
		if textGrid[row][searchCol] == targetChar {
			// Return position one after the target
			return row, searchCol + 1
		}
	}
	
	// Character not found in current line, return original position
	return row, col
}