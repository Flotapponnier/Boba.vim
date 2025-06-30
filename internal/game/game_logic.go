package game

import (
	"math/rand"
	"strings"
	"time"
)

// Game constants
const (
	MAX_LINE_LENGTH = 20
	INITIAL_PEARLS = 1
)

// Map values
const (
	EMPTY  = 0
	PLAYER = 1
	PEARL  = 3
)

// InitializeGameSession creates a new game with text grid and game map
func InitializeGameSession() map[string]interface{} {
	textGrid := createTextLines()
	gameMap := createGameMap(textGrid)
	
	return map[string]interface{}{
		"text_grid":        textGrid,
		"game_map":         gameMap,
		"player_pos":       map[string]int{"row": 0, "col": 0},
		"preferred_column": 0,
	}
}

// createTextLines creates text grid from game text (variable row lengths)
func createTextLines() [][]string {
	text := `Welcome to boba.vim !
This game is here to help you learn the vim motions fundamental,
it's a long journey over there ! but with patience,
determination !
Florent.`
	
	// Replace newlines with spaces and split into words
	text = strings.ReplaceAll(text, "\n", " ")
	text = strings.TrimSpace(text)
	words := strings.Fields(text)
	
	var lines []string
	currentLine := ""
	
	for _, word := range words {
		testLine := currentLine
		if currentLine != "" {
			testLine += " "
		}
		testLine += word
		
		if len(testLine) <= MAX_LINE_LENGTH {
			currentLine = testLine
		} else {
			if currentLine != "" {
				lines = append(lines, currentLine)
			}
			currentLine = word
		}
	}
	
	if currentLine != "" {
		lines = append(lines, currentLine)
	}
	
	// Convert lines to character grid
	var grid [][]string
	for _, line := range lines {
		row := make([]string, len(line))
		for i, char := range line {
			row[i] = string(char)
		}
		grid = append(grid, row)
	}
	
	return grid
}

// createGameMap creates initial game map with player at (0,0)
func createGameMap(textGrid [][]string) [][]int {
	gameMap := make([][]int, len(textGrid))
	
	for rowIdx, row := range textGrid {
		mapRow := make([]int, len(row))
		for colIdx := range row {
			if rowIdx == 0 && colIdx == 0 {
				mapRow[colIdx] = PLAYER // Player position
			} else {
				mapRow[colIdx] = EMPTY // Empty space
			}
		}
		gameMap[rowIdx] = mapRow
	}
	
	// Place one pearl randomly
	placeNewPearl(gameMap, 0, 0)
	return gameMap
}

// placeNewPearl places a new pearl at random empty position
func placeNewPearl(gameMap [][]int, playerRow, playerCol int) {
	rand.Seed(time.Now().UnixNano())
	
	// Find all empty positions
	var emptyPositions [][2]int
	for rowIdx := 0; rowIdx < len(gameMap); rowIdx++ {
		for colIdx := 0; colIdx < len(gameMap[rowIdx]); colIdx++ {
			if gameMap[rowIdx][colIdx] == EMPTY && !(rowIdx == playerRow && colIdx == playerCol) {
				emptyPositions = append(emptyPositions, [2]int{rowIdx, colIdx})
			}
		}
	}
	
	// Place pearl at random empty position
	if len(emptyPositions) > 0 {
		pos := emptyPositions[rand.Intn(len(emptyPositions))]
		gameMap[pos[0]][pos[1]] = PEARL
	}
}

// PlaceNewPearl is the exported version for external use
func PlaceNewPearl(gameMap [][]int, excludeRow, excludeCol int) {
	placeNewPearl(gameMap, excludeRow, excludeCol)
}

// IsValidPosition checks if a position is within bounds
func IsValidPosition(row, col int, gameMap [][]int) bool {
	if row < 0 || row >= len(gameMap) {
		return false
	}
	return col >= 0 && col < len(gameMap[row])
}