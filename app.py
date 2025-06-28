from flask import Flask, render_template, jsonify
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
I'm sure you gonna make it ! 
Best of luck,
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

    total_rows = len(text_grid)
    if total_rows > 1:
        pearl_row = random.randint(1, total_rows - 1)
        pearl_col = random.randint(0, len(text_grid[pearl_row]) - 1)
        game_map[pearl_row][pearl_col] = 3

    return game_map


@app.route("/")
def index():
    return render_template("index.html")


@app.route("/api/play")
def play():
    text_grid = create_text_lines()
    game_map = create_game_map(text_grid)

    return render_template("game.html", text_grid=text_grid, game_map=game_map)


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
