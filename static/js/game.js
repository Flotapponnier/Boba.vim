import { initializeGame, gameState } from "./game_modules/gameState.js";
import { initializeMovement } from "./game_modules/movement.js";
import { initializeTutorialMode } from "./game_modules/tutorial.js";
import { initializeChatHistory } from "./game_modules/chat.js";
import { initializeMapToggle } from "./game_modules/map.js";
import { initializeBackToMenuButton } from "./game_modules/navigation.js";
import { initializeResponsiveScaling } from "./game_modules/responsive_scaling.js";
import * as chatModule from "./game_modules/chat.js";
import * as tutorialModule from "./game_modules/tutorial.js";
import * as feedbackModule from "./game_modules/feedback.js";
import * as displayModule from "./game_modules/display.js";
import * as responsiveScaling from "./game_modules/responsive_scaling.js";

import * as CONSTANTS from "./game_modules/constants.js";

// Make modules available globally
window.gameState = gameState;
window.chatModule = chatModule;
window.tutorialModule = tutorialModule;
window.feedbackModule = feedbackModule;
window.displayModule = displayModule;
window.responsiveScaling = responsiveScaling;

// Make constants globally available
window.MOVEMENT_KEYS = CONSTANTS.MOVEMENT_KEYS;
window.MOVEMENT_MESSAGES = CONSTANTS.MOVEMENT_MESSAGES;
window.BLOCKED_MESSAGES = CONSTANTS.BLOCKED_MESSAGES;
window.VALID_MOVEMENT_KEYS = CONSTANTS.VALID_MOVEMENT_KEYS;
window.TUTORIAL_CONFIG = CONSTANTS.TUTORIAL_CONFIG;
window.TUTORIAL_COMMANDS = CONSTANTS.TUTORIAL_COMMANDS;
window.MAP_CONFIG = CONSTANTS.MAP_CONFIG;
window.FEEDBACK_CONFIG = CONSTANTS.FEEDBACK_CONFIG;
window.CHAT_CONFIG = CONSTANTS.CHAT_CONFIG;
window.API_ENDPOINTS = CONSTANTS.API_ENDPOINTS;
window.UI_SELECTORS = CONSTANTS.UI_SELECTORS;

document.addEventListener("DOMContentLoaded", function () {
  initializeGame();
  initializeBackToMenuButton();
  initializeMovement();
  initializeMapToggle();
  initializeTutorialMode();
  initializeChatHistory();
  
  // Initialize responsive scaling after everything else is set up
  setTimeout(() => {
    initializeResponsiveScaling();
  }, 100);
});
