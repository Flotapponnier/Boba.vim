"""
Grid Map Module for Boba.vim
Handles the creation and management of 2D grid maps from text
"""

from typing import List, Dict, Any, Optional, Tuple


class GridCell:
    """Represents a single cell in the grid"""

    def __init__(self, letter: str, row: int, col: int, is_empty: bool = False):
        self.letter = letter
        self.row = row
        self.col = col
        self.is_empty = is_empty
        self.is_highlighted = False
        self.metadata = {}

    def to_dict(self) -> Dict[str, Any]:
        """Convert cell to dictionary for JSON serialization"""
        return {
            "letter": self.letter,
            "row": self.row,
            "col": self.col,
            "is_empty": self.is_empty,
            "is_highlighted": self.is_highlighted,
            "metadata": self.metadata,
        }

    def __repr__(self):
        return f"GridCell('{self.letter}', {self.row}, {self.col})"


class GridMap:
    """
    Manages a 2D grid map created from text
    """

    def __init__(self, text: str, grid_width: int = 10):
        self.original_text = text
        self.processed_text = text.strip().upper()
        self.grid_width = grid_width
        self.grid_height = 0
        self.grid: List[List[GridCell]] = []
        self.cell_map: Dict[Tuple[int, int], GridCell] = {}

        # Generate the grid
        self._create_grid()

    def _create_grid(self) -> None:
        """Create the 2D grid from the processed text"""
        text_length = len(self.processed_text)
        self.grid_height = (text_length + self.grid_width - 1) // self.grid_width

        char_index = 0

        for row in range(self.grid_height):
            current_row = []
            for col in range(self.grid_width):
                if char_index < text_length:
                    letter = self.processed_text[char_index]
                    is_empty = False
                else:
                    letter = " "
                    is_empty = True

                cell = GridCell(letter, row, col, is_empty)
                current_row.append(cell)
                self.cell_map[(row, col)] = cell
                char_index += 1

            self.grid.append(current_row)

    def get_cell(self, row: int, col: int) -> Optional[GridCell]:
        """Get a specific cell by coordinates"""
        return self.cell_map.get((row, col))

    def get_letter_at(self, row: int, col: int) -> str:
        """Get the letter at specific coordinates"""
        cell = self.get_cell(row, col)
        return cell.letter if cell else " "

    def highlight_cell(self, row: int, col: int, highlight: bool = True) -> bool:
        """Highlight or unhighlight a specific cell"""
        cell = self.get_cell(row, col)
        if cell:
            cell.is_highlighted = highlight
            return True
        return False

    def highlight_text_sequence(self, start_row: int = 0, start_col: int = 0) -> int:
        """Highlight cells that contain actual text (non-empty)"""
        highlighted_count = 0

        for row in self.grid:
            for cell in row:
                if not cell.is_empty and cell.letter != " ":
                    cell.is_highlighted = True
                    highlighted_count += 1

        return highlighted_count

    def clear_all_highlights(self) -> None:
        """Remove highlights from all cells"""
        for row in self.grid:
            for cell in row:
                cell.is_highlighted = False

    def get_grid_as_lists(self) -> List[List[str]]:
        """Get the grid as a simple 2D list of letters (for template compatibility)"""
        return [[cell.letter for cell in row] for row in self.grid]

    def get_grid_as_dict_array(self) -> List[List[Dict[str, Any]]]:
        """Get the grid as a 2D list of cell dictionaries"""
        return [[cell.to_dict() for cell in row] for row in self.grid]

    def to_dict(self) -> Dict[str, Any]:
        """Convert the entire grid map to a dictionary for JSON serialization"""
        return {
            "original_text": self.original_text,
            "processed_text": self.processed_text,
            "grid_width": self.grid_width,
            "grid_height": self.grid_height,
            "grid": self.get_grid_as_dict_array(),
            "simple_grid": self.get_grid_as_lists(),
            "dimensions": {"width": self.grid_width, "height": self.grid_height},
            "stats": {
                "total_cells": self.grid_width * self.grid_height,
                "text_cells": len(
                    [c for row in self.grid for c in row if not c.is_empty]
                ),
                "empty_cells": len([c for row in self.grid for c in row if c.is_empty]),
            },
        }

    def get_text_coordinates(self) -> List[Tuple[int, int, str]]:
        """Get coordinates and letters of all non-empty cells"""
        coordinates = []
        for row in self.grid:
            for cell in row:
                if not cell.is_empty and cell.letter != " ":
                    coordinates.append((cell.row, cell.col, cell.letter))
        return coordinates

    def reconstruct_text(self) -> str:
        """Reconstruct the original text from the grid"""
        text_chars = []
        for row in self.grid:
            for cell in row:
                if not cell.is_empty:
                    text_chars.append(cell.letter)
        return "".join(text_chars).strip()

    def __repr__(self):
        return f"GridMap('{self.original_text}', {self.grid_width}x{self.grid_height})"

    def __str__(self):
        """String representation showing the grid visually"""
        lines = []
        for row in self.grid:
            line = " ".join(
                f"[{cell.letter}]" if not cell.is_empty else "[ ]" for cell in row
            )
            lines.append(line)
        return "\n".join(lines)


class GridMapFactory:
    """Factory class for creating different types of grid maps"""

    @staticmethod
    def create_welcome_map(grid_width: int = 8) -> GridMap:
        """Create the default welcome grid map"""
        return GridMap("Welcome to Boba.vim !", grid_width)

    @staticmethod
    def create_custom_map(text: str, grid_width: int = 10) -> GridMap:
        """Create a custom grid map with specified text"""
        return GridMap(text, grid_width)

    @staticmethod
    def create_square_map(text: str) -> GridMap:
        """Create a square grid map (width = height)"""
        text_length = len(text.strip())
        grid_size = int(text_length**0.5) + (1 if text_length**0.5 % 1 != 0 else 0)
        return GridMap(text, grid_size)

    @staticmethod
    def create_compact_map(text: str) -> GridMap:
        """Create a compact grid map with minimal empty cells"""
        text_length = len(text.strip())

        # Find the best width that minimizes empty cells
        best_width = text_length
        min_empty_cells = text_length

        for width in range(1, text_length + 1):
            height = (text_length + width - 1) // width
            total_cells = width * height
            empty_cells = total_cells - text_length

            if empty_cells < min_empty_cells:
                min_empty_cells = empty_cells
                best_width = width

        return GridMap(text, best_width)
