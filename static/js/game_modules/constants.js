// ================================
// MOVEMENT CONFIGURATION
// ================================
export const MOVEMENT_KEYS = {
  // Character & line movement
  h: { direction: "h", description: "LEFT ←" },
  j: { direction: "j", description: "DOWN ↓" },
  k: { direction: "k", description: "UP ↑" },
  l: { direction: "l", description: "RIGHT →" },

  // Word motions
  w: { direction: "w", description: "WORD FORWARD →" },
  W: { direction: "W", description: "WORD FORWARD (space-separated) →" },
  b: { direction: "b", description: "WORD BACK ←" },
  B: { direction: "B", description: "WORD BACK (space-separated) ←" },
  e: { direction: "e", description: "END WORD →" },
  E: { direction: "E", description: "END WORD (space-separated) →" },

  // Line motions
  0: { direction: "0", description: "LINE START ←" },
  $: { direction: "$", description: "LINE END →" },
  "^": { direction: "^", description: "FIRST NON-BLANK CHAR ←" },
  g_: { direction: "g_", description: "LAST NON-BLANK CHAR →" },

  // Screen positions
  H: { direction: "H", description: "TOP OF SCREEN" },
  M: { direction: "M", description: "MIDDLE OF SCREEN" },
  L: { direction: "L", description: "BOTTOM OF SCREEN" },

  // Paragraphs
  "{": { direction: "{", description: "PREV PARAGRAPH ←" },
  "}": { direction: "}", description: "NEXT PARAGRAPH →" },

  // Sentences
  "(": { direction: "(", description: "PREV SENTENCE ←" },
  ")": { direction: ")", description: "NEXT SENTENCE →" },

  // File
  gg: { direction: "gg", description: "TOP OF FILE" },
  G: { direction: "G", description: "BOTTOM OF FILE" },

  // Character search (line local)
  f: { direction: "f<char>", description: "FIND CHAR →" },
  F: { direction: "F<char>", description: "FIND CHAR ←" },
  t: { direction: "t<char>", description: "TILL CHAR →" },
  T: { direction: "T<char>", description: "TILL CHAR ←" },
  ";": { direction: ";", description: "REPEAT SEARCH" },
  ",": { direction: ",", description: "REVERSE REPEAT SEARCH" },

  // Match
  "%": { direction: "%", description: "MATCHING BRACKET" },
};

// Generate movement messages dynamically
export const MOVEMENT_MESSAGES = Object.fromEntries(
  Object.entries(MOVEMENT_KEYS).map(([key, config]) => [
    key,
    `You pressed ${key} to go ${config.description}`,
  ]),
);

export const BLOCKED_MESSAGES = Object.fromEntries(
  Object.entries(MOVEMENT_KEYS).map(([key, config]) => [
    key,
    `You pressed ${key} to go ${config.description} - BLOCKED!`,
  ]),
);

// Get list of valid movement keys
export const VALID_MOVEMENT_KEYS = Object.keys(MOVEMENT_KEYS);

// ================================
// TUTORIAL CONFIGURATION
// ================================
export const TUTORIAL_CONFIG = {
  TOGGLE_KEY: "+",
  COLORS: {
    ACTIVATED: "#9b59b6",
    INSTRUCTION: "#3498db",
    CORRECT: "#27ae60",
    WRONG: "#e74c3c",
  },
  TIMINGS: {
    ACTIVATION_DELAY: 1000,
    NEXT_COMMAND_DELAY: 1000,
    WRONG_ANSWER_DISPLAY: 1500,
  },
  MESSAGES: {
    ACTIVATED: "TUTORIAL MODE ACTIVATED",
    DEACTIVATED: "Tutorial mode deactivated",
    DEFAULT_WELCOME:
      "Welcome! Use HJKL to move | Press - to show map | Press + to activate tutorial",
  },
};

// Generate tutorial commands from movement configuration
export const TUTORIAL_COMMANDS = VALID_MOVEMENT_KEYS.map((key) => ({
  key: key,
  message: `Press ${key} to go ${MOVEMENT_KEYS[key].description}`,
}));

// ================================
// MAP CONFIGURATION
// ================================
export const MAP_CONFIG = {
  TOGGLE_KEY: "-",
  MESSAGES: {
    SHOWN: "Map SHOWN",
    HIDDEN: "Map HIDDEN",
    TOGGLE: "You pressed - to toggle map",
  },
  DISPLAY_ELEMENT_ID: "mapDisplay",
  FEEDBACK_TIMEOUT: 2000,
  FEEDBACK_COLOR: "#4ecdc4",
};

// ================================
// FEEDBACK CONFIGURATION
// ================================
export const FEEDBACK_CONFIG = {
  MOVEMENT_COLOR: "#ff6b6b",
  ANIMATION_DURATION: 300,
  ANIMATION_TYPE: "pulse 0.3s ease-in-out",
};

// ================================
// CHAT CONFIGURATION
// ================================
export const CHAT_CONFIG = {
  TOGGLE_KEY: "\\",
  MAX_HISTORY: 50,
  WINDOW_ID: "chatHistoryWindow",
};

// ================================
// GAME STATE
// ================================
export const GAME_STATE_DEFAULT = {
  TUTORIAL_MODE: false,
};

// ================================
// API ENDPOINTS
// ================================
export const API_ENDPOINTS = {
  MOVE: "/api/move",
  GAME_STATE: "/api/game-state",
  PLAY_TUTORIAL: "/api/playtutorial",
  PLAY_ONLINE: "/api/playonline",
  SET_USERNAME: "/api/set-username",
  LEADERBOARD: "/api/leaderboard",
  PLAYER_STATS: "/api/player-stats",
};

// ================================
// UI SELECTORS
// ================================
export const UI_SELECTORS = {
  HEADER_INFO: ".header-info",
  SCORE_ELEMENT: "#score",
  BACK_MENU_BUTTON: "#backMenu",
  GAME_KEYS: ".key",
  MAP_DISPLAY: "#mapDisplay",
  MAP_GRID: "#mapGrid",
};
