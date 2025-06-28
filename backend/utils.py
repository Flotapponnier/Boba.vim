"""
Utility functions for the Boba.vim backend
"""

from typing import List, Tuple, Dict, Any
import re


class TextProcessor:
    """Utility class for processing text before grid creation"""

    @staticmethod
    def clean_text(text: str, preserve_case: bool = False) -> str:
        """Clean and normalize text for grid processing"""
        # Remove extra whitespace
        cleaned = re.sub(r"\s+", " ", text.strip())

        if not preserve_case:
            cleaned = cleaned.upper()

        return cleaned

    @staticmethod
    def filter_printable(text: str) -> str:
        """Keep only printable ASCII characters"""
        return "".join(char for char in text if 32 <= ord(char) <= 126)

    @staticmethod
    def calculate_optimal_dimensions(text_length: int) -> Tuple[int, int]:
        """Calculate optimal grid dimensions for given text length"""
        if text_length == 0:
            return 1, 1

        # Try to find dimensions close to square
        width = int(text_length**0.5)

        # Adjust to minimize empty cells
        while width > 0:
            height = (text_length + width - 1) // width
            if width * height - text_length <= width:  # Acceptable waste
                return width, height
            width -= 1

        return text_length, 1  # Fallback to single row


class GridValidator:
    """Validation utilities for grid operations"""

    @staticmethod
    def validate_coordinates(
        row: int, col: int, grid_height: int, grid_width: int
    ) -> bool:
        """Check if coordinates are within grid bounds"""
        return 0 <= row < grid_height and 0 <= col < grid_width

    @staticmethod
    def validate_grid_dimensions(width: int, height: int) -> bool:
        """Validate grid dimensions are reasonable"""
        return 1 <= width <= 50 and 1 <= height <= 50

    @staticmethod
    def validate_text_length(text: str, max_length: int = 1000) -> bool:
        """Validate text length is reasonable for grid creation"""
        return 0 <= len(text) <= max_length


class GridExporter:
    """Export grid data in various formats"""

    @staticmethod
    def to_ascii_art(grid_map) -> str:
        """Convert grid to ASCII art representation"""
        lines = []
        lines.append("+" + "-" * (grid_map.grid_width * 3 - 1) + "+")

        for row in grid_map.grid:
            line = "|"
            for cell in row:
                if cell.letter == " ":
                    line += "   "
                else:
                    line += f" {cell.letter} "
            line = line.rstrip() + "|"
            lines.append(line)

        lines.append("+" + "-" * (grid_map.grid_width * 3 - 1) + "+")
        return "\n".join(lines)

    @staticmethod
    def to_csv(grid_map) -> str:
        """Convert grid to CSV format"""
        lines = []
        for row in grid_map.grid:
            csv_row = ",".join(f'"{cell.letter}"' for cell in row)
            lines.append(csv_row)
        return "\n".join(lines)

    @staticmethod
    def to_matrix_notation(grid_map) -> str:
        """Convert grid to mathematical matrix notation"""
        lines = ["["]
        for i, row in enumerate(grid_map.grid):
            row_str = " [" + ", ".join(f"'{cell.letter}'" for cell in row) + "]"
            if i < len(grid_map.grid) - 1:
                row_str += ","
            lines.append(row_str)
        lines.append("]")
        return "\n".join(lines)


class GameLogic:
    """Game-specific logic utilities"""

    @staticmethod
    def calculate_score(action: str, cell_letter: str, is_correct: bool = True) -> int:
        """Calculate score based on game action"""
        scores = {
            "click": 10,
            "highlight": 5,
            "correct_sequence": 50,
            "word_completion": 100,
        }

        base_score = scores.get(action, 0)

        # Bonus for non-space characters
        if cell_letter and cell_letter != " ":
            base_score += 5

        # Penalty for incorrect actions
        if not is_correct:
            base_score = max(0, base_score - 10)

        return base_score

    @staticmethod
    def check_word_completion(
        highlighted_cells: List[Tuple[int, int]], target_text: str, grid_map
    ) -> bool:
        """Check if highlighted cells form the target text"""
        # Sort cells by position (row-major order)
        sorted_cells = sorted(highlighted_cells)

        # Extract letters from highlighted cells
        highlighted_text = ""
        for row, col in sorted_cells:
            cell = grid_map.get_cell(row, col)
            if cell and not cell.is_empty:
                highlighted_text += cell.letter

        return highlighted_text.strip() == target_text.strip().upper()

    @staticmethod
    def get_next_hint_cell(current_progress: int, grid_map) -> Tuple[int, int]:
        """Get coordinates of the next cell to highlight as a hint"""
        text_coordinates = grid_map.get_text_coordinates()

        if current_progress < len(text_coordinates):
            row, col, letter = text_coordinates[current_progress]
            return row, col

        return -1, -1  # No more hints available


# Configuration constants
class Config:
    """Configuration constants for the backend"""

    DEFAULT_GRID_WIDTH = 8
    MAX_GRID_WIDTH = 20
    MAX_GRID_HEIGHT = 20
    MAX_TEXT_LENGTH = 500
    DEFAULT_GAME_TEXT = "Welcome to Boba.vim !"

    # Scoring
    SCORE_PER_CLICK = 10
    SCORE_PER_HIGHLIGHT = 5
    SCORE_WORD_COMPLETION = 100

    # Grid display
    EMPTY_CELL_CHAR = " "
    HIGHLIGHT_DURATION_MS = 1000
