import { initializePlayButton } from "./index_modules/playButton.js";
import { initializeLeaderboardButton } from "./index_modules/leaderboard.js";
import { initializeOnlineButton } from "./index_modules/onlineButton.js";
import { initializeUsernameInput } from "./index_modules/usernameInput.js";

document.addEventListener("DOMContentLoaded", function () {
  console.log("Initializing index page modules...");

  initializePlayButton();
  initializeOnlineButton();
  initializeLeaderboardButton();
  initializeUsernameInput();

  console.log("All index modules initialized successfully");
});
