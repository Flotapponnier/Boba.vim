import os
from dotenv import load_dotenv
from backend.game.list_movement import *

load_dotenv()


class Config:
    """Application configuration"""

    SECRET_KEY = os.getenv("SECRET_KEY", "boba-secret-key")
    DEBUG = True
    HOST = "127.0.0.1"
    PORT = 5000

    # Game configuration
    MAX_LINE_LENGTH = 20
    PEARL_POINTS = 100

    # Movement configuration - centralized key mapping
    MOVEMENT_FUNCTIONS = {
        "h": move_h,
        "j": move_j,
        "k": move_k,
        "l": move_l,
        "w": move_w,
        "b": move_b,
        "e": move_e,
        "0": move_0,
        "$": move_dollar,
    }

    # Get list of valid movement keys
    VALID_MOVEMENT_KEYS = list(MOVEMENT_FUNCTIONS.keys())

    # File monitoring for auto-reload
    EXTRA_FILES = [
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
    ]
