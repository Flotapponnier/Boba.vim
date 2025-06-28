import { initializeGame, gameState } from "./game_modules/gameState.js";
import { initializeMovement } from "./game_modules/movement.js";
import { initializeTutorialMode } from "./game_modules/tutorial.js";
import { initializeChatHistory } from "./game_modules/chat.js";
import { initializeMapToggle } from "./game_modules/map.js";
import { initializeBackToMenuButton } from "./game_modules/navigation.js";
import * as chatModule from "./game_modules/chat.js";
import * as tutorialModule from "./game_modules/tutorial.js";
import * as feedbackModule from "./game_modules/feedback.js";
import * as displayModule from "./game_modules/display.js";

window.gameState = gameState;
window.chatModule = chatModule;
window.tutorialModule = tutorialModule;
window.feedbackModule = feedbackModule;
window.displayModule = displayModule;

document.addEventListener("DOMContentLoaded", function () {
  initializeGame();
  initializeBackToMenuButton();
  initializeMovement();
  initializeMapToggle();
  initializeTutorialMode();
  initializeChatHistory();
});
