# ================================
# MOVEMENT HANDLER CLASS
# ================================


class MovementHandler:
    """Handles all player movement logic"""

    @staticmethod
    def calculate_new_position(
        direction, current_row, current_col, game_map, preferred_column
    ):
        """Calculate new position using config movement functions"""
        from backend.config import Config

        if direction not in Config.MOVEMENT_FUNCTIONS:
            raise ValueError(f"Invalid direction: {direction}")

        movement_func = Config.MOVEMENT_FUNCTIONS[direction]
        return movement_func(current_row, current_col, game_map, preferred_column)

    @staticmethod
    def is_valid_position(new_row, new_col, game_map):
        """Check if the new position is within bounds"""
        if new_row < 0 or new_row >= len(game_map):
            return False
        if new_col < 0 or new_col >= len(game_map[new_row]):
            return False
        return True

    @staticmethod
    def process_move(game_map, current_row, current_col, new_row, new_col, score):
        """Process the actual move and handle pearl collection"""
        from backend.config import Config

        target_value = game_map[new_row][new_col]
        pearl_collected = False

        if target_value == 3:
            pearl_collected = True
            score += Config.PEARL_POINTS
            from backend.game.game_logic import place_new_pearl

            place_new_pearl(game_map, new_row, new_col)

        # Update game map
        game_map[current_row][current_col] = 0
        game_map[new_row][new_col] = 1

        return {
            "pearl_collected": pearl_collected,
            "new_score": score,
            "updated_map": game_map,
        }

    @staticmethod
    def get_available_movements():
        """Get list of all available movement keys"""
        from backend.config import Config

        return Config.VALID_MOVEMENT_KEYS
