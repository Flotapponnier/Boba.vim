# ================================
# MOVEMENT FUNCTIONS FOR EACH KEY
# ================================


def move_h(current_row, current_col, game_map, preferred_column):
    """Move LEFT (h key) - vim style"""
    return {
        "new_row": current_row,
        "new_col": current_col - 1,
        "preferred_column": current_col - 1,
    }


def move_l(current_row, current_col, game_map, preferred_column):
    """Move RIGHT (l key) - vim style"""
    return {
        "new_row": current_row,
        "new_col": current_col + 1,
        "preferred_column": current_col + 1,
    }


def move_j(current_row, current_col, game_map, preferred_column):
    """Move DOWN (j key) - vim style with preferred column"""
    new_row = current_row + 1
    if preferred_column is None:
        preferred_column = current_col

    if new_row < len(game_map):
        target_row_length = len(game_map[new_row])
        if preferred_column >= target_row_length:
            new_col = target_row_length - 1
        else:
            new_col = preferred_column
    else:
        new_col = current_col

    return {
        "new_row": new_row,
        "new_col": new_col,
        "preferred_column": preferred_column,
    }


def move_k(current_row, current_col, game_map, preferred_column):
    """Move UP (k key) - vim style with preferred column"""
    new_row = current_row - 1
    if preferred_column is None:
        preferred_column = current_col

    if new_row >= 0:
        target_row_length = len(game_map[new_row])
        if preferred_column >= target_row_length:
            new_col = target_row_length - 1
        else:
            new_col = preferred_column
    else:
        new_col = current_col

    return {
        "new_row": new_row,
        "new_col": new_col,
        "preferred_column": preferred_column,
    }


# ================================
# FUTURE VIM MOVEMENTS (EASY TO ADD)
# ================================


def move_w(current_row, current_col, game_map, preferred_column):
    """Move to next WORD (w key) - vim style"""
    # TODO: Implement word-forward movement
    return move_l(current_row, current_col, game_map, preferred_column)


def move_b(current_row, current_col, game_map, preferred_column):
    """Move to previous WORD (b key) - vim style"""
    # TODO: Implement word-backward movement
    return move_h(current_row, current_col, game_map, preferred_column)


def move_e(current_row, current_col, game_map, preferred_column):
    """Move to END of word (e key) - vim style"""
    # TODO: Implement end-of-word movement
    return move_l(current_row, current_col, game_map, preferred_column)


def move_0(current_row, current_col, game_map, preferred_column):
    """Move to START of line (0 key) - vim style"""
    return {
        "new_row": current_row,
        "new_col": 0,
        "preferred_column": 0,
    }


def move_dollar(current_row, current_col, game_map, preferred_column):
    """Move to END of line ($ key) - vim style"""
    if current_row < len(game_map):
        end_col = len(game_map[current_row]) - 1
    else:
        end_col = current_col

    return {
        "new_row": current_row,
        "new_col": end_col,
        "preferred_column": end_col,
    }
