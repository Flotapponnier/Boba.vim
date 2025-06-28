from flask import render_template, jsonify
from backend.controllers.game_controller import GameController


def register_routes(app):
    """Register all application routes"""

    @app.route("/")
    def index():
        """Home page"""
        return render_template("index.html")

    @app.route("/api/play")
    def play():
        """Start a new game"""
        try:
            game_data = GameController.initialize_game()
            return render_template(
                "game.html",
                text_grid=game_data["text_grid"],
                game_map=game_data["game_map"],
            )
        except Exception as e:
            return jsonify({"success": False, "error": str(e)}), 500

    @app.route("/api/move", methods=["POST"])
    def move_player():
        """Handle player movement"""
        return GameController.move_player()

    @app.route("/api/movements")
    def get_available_movements():
        """Get list of all available movement keys"""
        return GameController.get_available_movements()

    @app.route("/api/game-state")
    def get_game_state():
        """Get current game state"""
        return GameController.get_game_state()

    @app.route("/api/playtutorial")
    def play_tutorial():
        """Tutorial mode endpoint (to be implemented)"""
        return jsonify(
            {"message": "Tutorial will be implemented soon", "status": "not ready"}
        )

    @app.route("/api/playonline")
    def play_online():
        """Online mode endpoint (to be implemented)"""
        return jsonify(
            {"message": "Online game will be implemented soon", "status": "not ready"}
        )

    @app.errorhandler(404)
    def page_not_found(e):
        """404 error handler"""
        return render_template("404.html"), 404

    @app.errorhandler(500)
    def internal_server_error(e):
        """500 error handler"""
        return render_template("500.html"), 500
