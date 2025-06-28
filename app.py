from flask import Flask
from backend.config import Config
from backend.routes import register_routes


def create_app():
    """Application factory pattern"""
    app = Flask(__name__)
    app.config["SECRET_KEY"] = Config.SECRET_KEY

    register_routes(app)

    return app


def main():
    """Main application entry point"""
    app = create_app()

    app.run(
        debug=Config.DEBUG,
        host=Config.HOST,
        port=Config.PORT,
        extra_files=Config.EXTRA_FILES,
    )


if __name__ == "__main__":
    main()
