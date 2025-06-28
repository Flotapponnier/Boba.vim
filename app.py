from flask import Flask, render_template, jsonify
import os
from dotenv import load_dotenv

load_dotenv()

app = Flask(__name__)
app.config["SECRET_KEY"] = os.getenv("SECRET_KEY", "boba-secret-key")


def create_text_lines():
    """Convert text to lines with natural word breaks"""
    text = """Welcome to 
boba.vim !
This game 
is here to
help you learn 
the 
vim motions 
fundamental,
it's a 
long journey over 
there !
but with patience,
determination !
I'm sure you 
gonna make it ! 
Best of luck,
Florent."""

    lines = text.strip().split("\n")

    grid = []
    for line in lines:
        if line.strip():
            char_list = []
            for char in line:
                char_list.append(char)
            grid.append(char_list)

    return grid


@app.route("/")
def index():
    return render_template("index.html")


@app.route("/api/play")
def play():
    grid = create_text_lines()

    # Debug: print the grid structure
    print("Grid structure:")
    for i, row in enumerate(grid):
        print(f"Row {i}: {len(row)} characters - {''.join(row)}")

    return render_template("game.html", grid=grid)


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
            # Templates
            "templates/base.html",
            "templates/index.html",
            "templates/game.html",
            "templates/404.html",
            "templates/500.html",
            # CSS files
            "static/css/global.css",
            "static/css/index.css",
            "static/css/game.css",
            # JavaScript files
            "static/js/index.js",
            "static/js/game.js",
        ],
    )
