// ================================
// MOVEMENT CONFIGURATION
// ================================
export const MOVEMENT_KEYS = {
  // Basic vim movements
  h: { direction: "h", description: "LEFT ←" },
  j: { direction: "j", description: "DOWN ↓" },
  k: { direction: "k", description: "UP ↑" },
  l: { direction: "l", description: "RIGHT →" },

  // Future keys can be added here
  // w: { direction: "w", description: "WORD FORWARD →" },
  // b: { direction: "b", description: "WORD BACK ←" },
  // e: { direction: "e", description: "END WORD →" },
  // 0: { direction: "0", description: "LINE START ←" },
  // $: { direction: "$", description: "LINE END →" },
};

// Generate movement messages dynamically
export const MOVEMENT_MESSAGES = Object.fromEntries(
  Object.entries(MOVEMENT_KEYS).map(([key, config]) => [
    key,
    `You pressed ${key.toUpperCase()} to go ${config.description}`,
  ]),
);

export const BLOCKED_MESSAGES = Object.fromEntries(
  Object.entries(MOVEMENT_KEYS).map(([key, config]) => [
    key,
    `You pressed ${key.toUpperCase()} to go ${config.description} - BLOCKED!`,
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
  message: `Press ${key.toUpperCase()} to go ${MOVEMENT_KEYS[key].description}`,
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
