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
    return jsonify({"message": "Welcome to boba.vim!", "status": "ready"})


@app.route("/api/playtutorial")
def play_tutorial():
    return jsonify(
        {"message": " Tutorial will be implement soon", "status": "not ready"}
    )


@app.route("/api/playonline")
def play_online():
    return jsonify(
        {"message": " Online game will be implement soon", "status": "not ready"}
    )


if __name__ == "__main__":
    # Enable live reloading for all files
    app.run(
        debug=True,
        host="127.0.0.1",
        port=5000,
        extra_files=["templates/index.html", "static/style.css", "static/script.js"],
    )
