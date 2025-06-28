from flask import Flask, render_template, jsonify, request, session
import os
import random
from dotenv import load_dotenv

load_dotenv()

app = Flask(__name__)
app.config["SECRET_KEY"] = os.getenv("SECRET_KEY", "boba-secret-key")


def create_text_lines():
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
        if len(test_line) <= 20:
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
    game_map = []

    for row_idx, row in enumerate(text_grid):
        map_row = []
        for col_idx, char in enumerate(row):
            if row_idx == 0 and col_idx == 0:
                map_row.append(1)
            else:
                map_row.append(0)
        game_map.append(map_row)

    place_new_pearl(game_map, 0, 0)

    return game_map


@app.route("/")
def index():
    return render_template("index.html")


@app.route("/api/play")
def play():
    text_grid = create_text_lines()
    game_map = create_game_map(text_grid)

    session["text_grid"] = text_grid
    session["game_map"] = game_map
    session["player_pos"] = {"row": 0, "col": 0}
    session["score"] = 0

    return render_template("game.html", text_grid=text_grid, game_map=game_map)


@app.route("/api/move", methods=["POST"])
def move_player():
    data = request.get_json()
    direction = data.get("direction")

    game_map = session.get("game_map", [])
    player_pos = session.get("player_pos", {"row": 0, "col": 0})
    score = session.get("score", 0)

    current_row = player_pos["row"]
    current_col = player_pos["col"]

    if direction == "h":
        new_row, new_col = current_row, current_col - 1
    elif direction == "j":
        new_row, new_col = current_row + 1, current_col
    elif direction == "k":
        new_row, new_col = current_row - 1, current_col
    elif direction == "l":
        new_row, new_col = current_row, current_col + 1
    else:
        return jsonify({"success": False, "error": "Invalid direction"})

    if new_row < 0 or new_row >= len(game_map):
        return jsonify({"success": False, "error": "Out of bounds"})

    if new_col < 0 or new_col >= len(game_map[new_row]):
        return jsonify({"success": False, "error": "Out of bounds"})

    target_value = game_map[new_row][new_col]
    pearl_collected = False

    if target_value == 3:
        pearl_collected = True
        score += 100

        place_new_pearl(game_map, new_row, new_col)

    game_map[current_row][current_col] = 0
    game_map[new_row][new_col] = 1

    new_player_pos = {"row": new_row, "col": new_col}

    session["game_map"] = game_map
    session["player_pos"] = new_player_pos
    session["score"] = score

    return jsonify(
        {
            "success": True,
            "game_map": game_map,
            "player_pos": new_player_pos,
            "score": score,
            "pearl_collected": pearl_collected,
        }
    )


def place_new_pearl(game_map, player_row, player_col):
    empty_positions = []

    for row_idx in range(len(game_map)):
        for col_idx in range(len(game_map[row_idx])):
            if game_map[row_idx][col_idx] == 0 and not (
                row_idx == player_row and col_idx == player_col
            ):
                empty_positions.append((row_idx, col_idx))

    if empty_positions:
        pearl_row, pearl_col = random.choice(empty_positions)
        game_map[pearl_row][pearl_col] = 3


@app.route("/api/game-state")
def get_game_state():
    return jsonify(
        {
            "game_map": session.get("game_map", []),
            "player_pos": session.get("player_pos", {"row": 0, "col": 0}),
            "score": session.get("score", 0),
        }
    )


@app.route("/api/playtutorial")
def play_tutorial():
    return jsonify(
        {"message": "Tutorial will be implemented soon", "status": "not ready"}
    )


@app.route("/api/playonline")
def play_online():
    return jsonify(
        {"message": "Online game will be implemented soon", "status": "not ready"}
    )


@app.errorhandler(404)
def page_not_found(e):
    return render_template("404.html"), 404


@app.errorhandler(500)
def internal_server_error(e):
    return render_template("500.html"), 500


if __name__ == "__main__":
    app.run(
        debug=True,
        host="127.0.0.1",
        port=5000,
        extra_files=[
            "templates/base.html",
            "templates/index.html",
            "templates/game.html",
            "templates/404.html",
            "templates/500.html",
            "static/css/global.css",
            "static/css/index.css",
            "static/css/game.css",
            "static/js/index.js",
            "static/js/game.js",
        ],
    )
