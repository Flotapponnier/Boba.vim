package game

import (
	"math/rand"
	"strings"
	"time"
)

// Game constants
const (
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
	// Text patterns for vim practice with reasonable sizes (max ~50 chars width)
	textPatterns := []string{
		// Pattern 1: Welcome message
		`Welcome to boba.vim !
This game helps you learn vim motions fundamentals,
it's a long journey but with patience,
determination you'll master it!
Florent.`,
		
		// Pattern 2: JavaScript function
		`function movePlayer(direction) {
    if (direction === "up") {
        player.y -= 1;
    } else if (direction === "down") {
        player.y += 1;
    } else if (direction === "left") {
        player.x -= 1;
    } else if (direction === "right") {
        player.x += 1;
    }
    return player;
}`,

		// Pattern 3: Configuration syntax
		`server:
  host: localhost
  port: 8080
  ssl: true
  max_connections: 1000
database:
  url: postgres://localhost/boba_vim
  pool_size: 10
  timeout: 30s
cache:
  redis_url: redis://localhost:6379
  ttl: 3600`,

		// Pattern 4: Markdown guide
		`# Vim Motions Guide
## Basic Movement
- h: move left
- j: move down  
- k: move up
- l: move right

## Word Movement
- w: next word beginning
- b: previous word beginning
- e: end of word

## Line Movement
- 0: start of line
- $: end of line
- ^: first non-blank character
- g_: last non-blank character

## File Movement
- gg: top of file
- G: bottom of file`,

		// Pattern 5: JSON configuration
		`{
  "name": "boba-vim",
  "version": "1.0.0",
  "description": "Learn vim with boba tea!",
  "config": {
    "movement": ["h", "j", "k", "l"],
    "search": ["f", "F", "t", "T"],
    "navigation": ["w", "b", "e", "0", "$"]
  },
  "features": {
    "character_search": true,
    "word_movement": true,
    "line_navigation": true
  }
}`,

		// Pattern 6: CSS styles
		`.vim-game {
  background: #2c3e50;
  color: #ecf0f1;
  font-family: monospace;
  padding: 20px;
}

.player {
  position: absolute;
  width: 20px;
  height: 20px;
  background: #f39c12;
  border-radius: 50%;
  transition: all 0.2s ease;
}

.key {
  display: inline-block;
  padding: 8px 12px;
  margin: 2px;
  border: 2px solid #bdc3c7;
  border-radius: 4px;
}`,

		// Pattern 7: Python code
		`import random
import time

class BobaGame:
    def __init__(self):
        self.player_pos = {"x": 0, "y": 0}
        self.score = 0
        self.pearls = []
        
    def move_player(self, direction):
        if direction == "h":
            self.player_pos["x"] -= 1
        elif direction == "j":
            self.player_pos["y"] += 1
        elif direction == "k":
            self.player_pos["y"] -= 1
        elif direction == "l":
            self.player_pos["x"] += 1
            
    def collect_pearl(self):
        self.score += 100
        self.spawn_new_pearl()
        
    def spawn_new_pearl(self):
        x = random.randint(0, 20)
        y = random.randint(0, 15)
        self.pearls.append({"x": x, "y": y})`,

		// Pattern 8: Heavy spacing for ^ and g_ practice
		`        function calculateScore() {        
            let base = 1000;        
            let penalty = time * 10;        
                
            if (moves < 50) {        
                bonus = 200;        
            } else {        
                bonus = 0;        
            }        
            return base + bonus - penalty;        
        }`,

		// Pattern 9: SQL queries
		`SELECT u.username, u.email, p.score, p.completion_time
FROM users u
JOIN player_stats p ON u.id = p.user_id
WHERE p.score > 1000
  AND p.completion_time < 300
ORDER BY p.score DESC, p.completion_time ASC
LIMIT 10;

UPDATE game_sessions 
SET is_completed = true,
    final_score = current_score,
    end_time = NOW()
WHERE session_token = ? AND is_active = true;`,

		// Pattern 10: Spaced configuration for practice
		`     server:     
       host: localhost      
       port: 8080    
       ssl: true     
     database:    
       url: postgres://localhost/db     
       pool_size: 10      
     cache:    
       redis: localhost:6379     
       ttl: 3600      `,

		// Pattern 11: Mixed spacing challenge
		`         x = 1;           
              y = 2;        
      a = 4;            
             b = 5;      
     final = x + y + a + b;        `,

		// Pattern 12: Go code
		`package main

import (
    "fmt"
    "net/http"
    "log"
)

type Player struct {
    ID       int    ` + "`json:\"id\"`" + `
    Username string ` + "`json:\"username\"`" + `
    Score    int    ` + "`json:\"score\"`" + `
    Position struct {
        X int ` + "`json:\"x\"`" + `
        Y int ` + "`json:\"y\"`" + `
    } ` + "`json:\"position\"`" + `
}

func (p *Player) Move(direction string) {
    switch direction {
    case "h":
        p.Position.X--
    case "j":
        p.Position.Y++
    case "k":
        p.Position.Y--
    case "l":
        p.Position.X++
    }
}`,
	}
	
	// Randomly select one pattern
	rand.Seed(time.Now().UnixNano())
	text := textPatterns[rand.Intn(len(textPatterns))]
	
	// Split text into lines preserving all whitespace structure
	lines := strings.Split(text, "\n")
	
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