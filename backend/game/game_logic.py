import random
from backend.config import Config


def create_text_lines():
    """Create text grid from game text"""
    text = """Welcome to boba.vim !
This game is here to help you learn the vim motions fundamental,
it's a long journey over there ! but with patience,
determination !
Florent."""

    text = text.replace("\n", " ").strip()

    lines = []
    words = text.split()
    current_line = ""

    for word in words:
        test_line = current_line + (" " if current_line else "") + word
        if len(test_line) <= Config.MAX_LINE_LENGTH:
            current_line = test_line
        else:
            if current_line:
                lines.append(current_line)
            current_line = word

    if current_line:
        lines.append(current_line)

    grid = []
    for line in lines:
        grid.append(list(line))

    return grid


def create_game_map(text_grid):
    """Create initial game map with player at (0,0)"""
    game_map = []

    for row_idx, row in enumerate(text_grid):
        map_row = []
        for col_idx, char in enumerate(row):
            if row_idx == 0 and col_idx == 0:
                map_row.append(1)  # Player position
            else:
                map_row.append(0)  # Empty space
        game_map.append(map_row)

    place_new_pearl(game_map, 0, 0)
    return game_map


def place_new_pearl(game_map, player_row, player_col):
    """Place a new pearl at random empty position"""
    empty_positions = []

    for row_idx in range(len(game_map)):
        for col_idx in range(len(game_map[row_idx])):
            if game_map[row_idx][col_idx] == 0 and not (
                row_idx == player_row and col_idx == player_col
            ):
                empty_positions.append((row_idx, col_idx))

    if empty_positions:
        pearl_row, pearl_col = random.choice(empty_positions)
        game_map[pearl_row][pearl_col] = 3  # Pearl


def initialize_game_session():
    """Initialize a new game session"""
    text_grid = create_text_lines()
    game_map = create_game_map(text_grid)

    return {
        "text_grid": text_grid,
        "game_map": game_map,
        "player_pos": {"row": 0, "col": 0},
        "score": 0,
        "preferred_column": None,
    }
