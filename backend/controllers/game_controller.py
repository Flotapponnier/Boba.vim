from flask import jsonify, request, session
from backend.game.game_logic import initialize_game_session
from backend.game.movement import MovementHandler


class GameController:
    """Handles all game-related HTTP requests"""

    @staticmethod
    def move_player():
        """Handle player movement"""
        try:
            data = request.get_json()
            direction = data.get("direction")

            if not direction:
                return jsonify({"success": False, "error": "No direction provided"})

            game_state = GameController._get_game_state_from_session()

            try:
                movement_result = MovementHandler.calculate_new_position(
                    direction,
                    game_state["current_row"],
                    game_state["current_col"],
                    game_state["game_map"],
                    game_state["preferred_column"],
                )
            except ValueError as e:
                return jsonify({"success": False, "error": str(e)})

            new_row = movement_result["new_row"]
            new_col = movement_result["new_col"]
            preferred_column = movement_result["preferred_column"]

            if not MovementHandler.is_valid_position(
                new_row, new_col, game_state["game_map"]
            ):
                return jsonify({"success": False, "error": "Out of bounds"})

            move_result = MovementHandler.process_move(
                game_state["game_map"],
                game_state["current_row"],
                game_state["current_col"],
                new_row,
                new_col,
                game_state["score"],
            )

            GameController._update_session_with_move_result(
                move_result, new_row, new_col, preferred_column
            )

            return GameController._create_move_response(move_result, new_row, new_col)

        except Exception as e:
            return jsonify({"success": False, "error": f"Server error: {str(e)}"})

    @staticmethod
    def initialize_game():
        """Initialize a new game session"""
        try:
            game_data = initialize_game_session()

            session["text_grid"] = game_data["text_grid"]
            session["game_map"] = game_data["game_map"]
            session["player_pos"] = game_data["player_pos"]
            session["score"] = game_data["score"]
            session["preferred_column"] = game_data["preferred_column"]

            return game_data

        except Exception as e:
            raise Exception(f"Failed to initialize game: {str(e)}")

    @staticmethod
    def get_game_state():
        """Get current game state"""
        return jsonify(
            {
                "game_map": session.get("game_map", []),
                "player_pos": session.get("player_pos", {"row": 0, "col": 0}),
                "score": session.get("score", 0),
            }
        )

    @staticmethod
    def get_available_movements():
        """Get list of all available movement keys"""
        return jsonify(
            {
                "movements": MovementHandler.get_available_movements(),
                "total": len(MovementHandler.get_available_movements()),
            }
        )

    # ================================
    # PRIVATE HELPER METHODS
    # ================================

    @staticmethod
    def _get_game_state_from_session():
        """Extract game state from session"""
        game_map = session.get("game_map", [])
        player_pos = session.get("player_pos", {"row": 0, "col": 0})
        score = session.get("score", 0)
        preferred_column = session.get("preferred_column", None)

        return {
            "game_map": game_map,
            "current_row": player_pos["row"],
            "current_col": player_pos["col"],
            "score": score,
            "preferred_column": preferred_column,
        }

    @staticmethod
    def _update_session_with_move_result(
        move_result, new_row, new_col, preferred_column
    ):
        """Update session with move results"""
        new_player_pos = {"row": new_row, "col": new_col}
        session["game_map"] = move_result["updated_map"]
        session["player_pos"] = new_player_pos
        session["score"] = move_result["new_score"]
        session["preferred_column"] = preferred_column

    @staticmethod
    def _create_move_response(move_result, new_row, new_col):
        """Create JSON response for successful move"""
        return jsonify(
            {
                "success": True,
                "game_map": move_result["updated_map"],
                "player_pos": {"row": new_row, "col": new_col},
                "score": move_result["new_score"],
                "pearl_collected": move_result["pearl_collected"],
            }
        )
