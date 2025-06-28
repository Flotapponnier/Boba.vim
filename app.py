from flask import Flask, render_template, jsonify
import os
from dotenv import load_dotenv

load_dotenv()

app = Flask(__name__)
app.config["SECRET_KEY"] = os.getenv("SECRET_KEY", "boba-secret-key")


@app.route("/")
def index():
    return render_template("index.html")


@app.route("/api/play")
def play():
    return render_template("game.html")


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
