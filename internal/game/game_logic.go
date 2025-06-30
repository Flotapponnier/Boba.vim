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
	// Randomized text patterns for vim practice
	textPatterns := []string{
		// Pattern 1: Welcome message
		`Welcome to boba.vim !
This game is here to help you learn the vim motions fundamental,
it's a long journey over there ! but with patience,
determination !
Florent.`,
		
		// Pattern 2: Regular JavaScript code
		`function movePlayer(direction) {
    if (direction === "up") {
        player.y -= 1;
    } else if (direction === "down") {
        player.y += 1;
    }
    return player;
}`,

		// Pattern 3: Configuration syntax
		`server:
  host: localhost
  port: 8080
  ssl: true
database:
  url: postgres://localhost
  pool_size: 10`,

		// Pattern 4: Markdown syntax
		`# Vim Motions Guide
## Basic Movement
- h: left
- j: down  
- k: up
- l: right
## Word Movement
Use w, b, e for words!`,

		// Pattern 5: JSON structure
		`{
  "name": "boba-vim",
  "version": "1.0.0",
  "config": {
    "movement": ["h", "j", "k", "l"],
    "search": ["f", "F", "t", "T"]
  }
}`,

		// Pattern 6: CSS styles
		`.vim-game {
  background: #2c3e50;
  color: #ecf0f1;
  font-family: monospace;
}

.player {
  position: absolute;
  width: 20px;
  height: 20px;
}`,

		// Pattern 7: Code with heavy indentation (^ and g_ practice)
		`        function calculate() {    
            let base = 1000;     
            let penalty = time * 10;      
                
            if (moves < 50) {     
                bonus = 200;     
            } else {      
                bonus = 0;     
            }     
            return base + bonus;     
        }`,

		// Pattern 8: Spaced configuration (^ and g_ practice)
		`     server:     
       host: localhost      
       port: 8080    
     database:    
       url: postgres://localhost     
       pool_size: 10      `,

		// Pattern 9: Mixed spacing challenge (^ and g_ practice)
		`         x = 1;           
              y = 2;        
      a = 4;            
             b = 5;      
     final = end;        `,
	}
	
	// Randomly select one pattern
	rand.Seed(time.Now().UnixNano())
	text := textPatterns[rand.Intn(len(textPatterns))]
	
	// Split text into lines while preserving whitespace structure
	rawLines := strings.Split(text, "\n")
	
	var lines []string
	for _, line := range rawLines {
		// Preserve leading and trailing spaces, but limit line length
		if len(line) <= MAX_LINE_LENGTH {
			lines = append(lines, line)
		} else {
			// If line is too long, break it while preserving leading spaces
			leadingSpaces := ""
			trimmed := strings.TrimLeft(line, " \t")
			if len(trimmed) < len(line) {
				leadingSpaces = line[:len(line)-len(trimmed)]
			}
			
			words := strings.Fields(trimmed)
			currentLine := leadingSpaces
			
			for _, word := range words {
				testLine := currentLine
				if strings.TrimSpace(currentLine) != "" {
					testLine += " "
				}
				testLine += word
				
				if len(testLine) <= MAX_LINE_LENGTH {
					currentLine = testLine
				} else {
					if strings.TrimSpace(currentLine) != "" {
						lines = append(lines, currentLine)
					}
					currentLine = leadingSpaces + word
				}
			}
			
			if strings.TrimSpace(currentLine) != "" {
				lines = append(lines, currentLine)
			}
		}
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