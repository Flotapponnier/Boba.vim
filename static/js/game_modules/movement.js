export function initializeMovement() {
  let lastKeyPressed = null;
  let keyReleased = true;

  document.addEventListener("keydown", function (event) {
    const key = event.key;

    if (window.VALID_MOVEMENT_KEYS.includes(key)) {
      if (event.repeat || (lastKeyPressed === key && !keyReleased)) {
        event.preventDefault();
        return;
      }

      lastKeyPressed = key;
      keyReleased = false;

      if (window.gameState.tutorialMode) {
        window.tutorialModule.handleTutorialMovement(key);
      } else {
        window.feedbackModule.showMovementFeedback(key);
      }

      movePlayer(key);
      event.preventDefault();
    }
  });

  document.addEventListener("keyup", function (event) {
    const key = event.key;

    if (window.VALID_MOVEMENT_KEYS.includes(key)) {
      if (lastKeyPressed === key) {
        keyReleased = true;
      }
    }
  });
}

let movePending = false;
let lastMoveTime = 0;
const MOVE_COOLDOWN = 20;

export async function movePlayer(direction) {
  if (window.gameCompleted) {
    console.log("Game is completed, movement disabled");
    return;
  }

  if (movePending) {
    console.log("Move already pending, ignoring");
    return;
  }

  const now = Date.now();
  if (now - lastMoveTime < MOVE_COOLDOWN) {
    console.log("Move too fast, ignoring");
    return;
  }

  movePending = true;
  lastMoveTime = now;

  try {
    const response = await fetch(window.API_ENDPOINTS.MOVE, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        direction: direction,
      }),
    });

    const result = await response.json();

    if (result.success) {
      handleSuccessfulMove(result, direction);
    } else {
      handleBlockedMove(result, direction);
    }
  } catch (error) {
    console.error("Error moving player:", error);
  } finally {
    movePending = false;
  }
}

function handleSuccessfulMove(result, direction) {
  window.displayModule.updateGameDisplay(result.game_map);
  window.displayModule.updateScore(result.score);

  if (result.pearl_collected) {
    handlePearlCollection(direction);
  }
  if (result.is_completed) {
    handleGameCompletion(result);
  }
}

function handlePearlCollection(direction) {
  console.log("Pearl collected! +100 points");
  window.chatModule.addToChatHistory("🧋 PEARL COLLECTED! +100 points!");

  if (!window.gameState.tutorialMode) {
    showPearlCollectionFeedback(direction);
  }
}

function showPearlCollectionFeedback(direction) {
  const headerInfo = document.querySelector(window.UI_SELECTORS.HEADER_INFO);
  if (!headerInfo) return;

  if (!headerInfo.dataset.originalContent) {
    headerInfo.dataset.originalContent = headerInfo.innerHTML;
  }

  headerInfo.innerHTML = `<strong style="color: #ffd700; animation: pulse 0.5s ease-in-out;">🧋 PEARL COLLECTED! +100 points!</strong>`;

  setTimeout(() => {
    const message =
      window.MOVEMENT_MESSAGES[direction] ||
      `You pressed ${direction.toUpperCase()}`;
    headerInfo.innerHTML = `<strong style="color: ${window.FEEDBACK_CONFIG.MOVEMENT_COLOR};">${message}</strong>`;
  }, 1500);
}

function handleBlockedMove(result, direction) {
  console.log("Move blocked:", result.error);

  const message =
    window.BLOCKED_MESSAGES[direction] ||
    `You pressed ${direction.toUpperCase()} - BLOCKED!`;
  window.chatModule.addToChatHistory(message);

  if (!window.gameState.tutorialMode) {
    showBlockedMoveFeedback(message);
  }
}

function showBlockedMoveFeedback(message) {
  const headerInfo = document.querySelector(window.UI_SELECTORS.HEADER_INFO);
  if (!headerInfo) return;

  if (!headerInfo.dataset.originalContent) {
    headerInfo.dataset.originalContent = headerInfo.innerHTML;
  }

  headerInfo.innerHTML = `<strong style="color: #e74c3c;">${message}</strong>`;
}

function handleGameCompletion(result) {
  console.log("Game completed! Final score:", result.final_score);

  // Disable further movement
  window.gameCompleted = true;

  // Show completion message
  const headerInfo = document.querySelector(window.UI_SELECTORS.HEADER_INFO);
  if (headerInfo) {
    const timeText = result.completion_time
      ? `Time: ${formatTime(result.completion_time)}`
      : "Time: --:--";

    headerInfo.innerHTML = `<strong style="color: #ffd700; font-size: 1.2em; animation: pulse 1s ease-in-out infinite;">
      🎉 GAME COMPLETED! 🎉<br>
      Final Score: ${result.final_score} | ${timeText}
    </strong>`;
  }

  window.chatModule.addToChatHistory(
    `🎉 CONGRATULATIONS! Game completed with score ${result.final_score}!`,
  );

  setTimeout(() => {
    showCompletionModal(result);
  }, 2000);
}

function formatTime(seconds) {
  const minutes = Math.floor(seconds / 60);
  const remainingSeconds = seconds % 60;
  return `${minutes}:${remainingSeconds.toString().padStart(2, "0")}`;
}

function showCompletionModal(result) {
  // Create modal HTML
  const modalHTML = `
    <div id="completionModal" style="
      position: fixed;
      top: 0;
      left: 0;
      width: 100%;
      height: 100%;
      background: rgba(0, 0, 0, 0.8);
      display: flex;
      justify-content: center;
      align-items: center;
      z-index: 1000;
    ">
      <div style="
        background: #2c3e50;
        color: white;
        padding: 2rem;
        border-radius: 10px;
        text-align: center;
        max-width: 500px;
        width: 90%;
        box-shadow: 0 0 20px rgba(255, 215, 0, 0.5);
      ">
        <h2 style="color: #ffd700; margin-bottom: 1rem;">🎉 Game Completed! 🎉</h2>
        <div style="margin: 1rem 0;">
          <p><strong>Final Score:</strong> ${result.final_score}</p>
          <p><strong>Completion Time:</strong> ${result.completion_time ? formatTime(result.completion_time) : "--:--"}</p>
        </div>
        <div style="margin: 1.5rem 0;">
          <button id="viewLeaderboard" style="
            background: #3498db;
            color: white;
            border: none;
            padding: 0.8rem 1.5rem;
            margin: 0.5rem;
            border-radius: 5px;
            cursor: pointer;
            font-size: 1rem;
          ">View Leaderboard</button>
          <button id="playAgain" style="
            background: #27ae60;
            color: white;
            border: none;
            padding: 0.8rem 1.5rem;
            margin: 0.5rem;
            border-radius: 5px;
            cursor: pointer;
            font-size: 1rem;
          ">Play Again</button>
          <button id="backToMenu" style="
            background: #95a5a6;
            color: white;
            border: none;
            padding: 0.8rem 1.5rem;
            margin: 0.5rem;
            border-radius: 5px;
            cursor: pointer;
            font-size: 1rem;
          ">Back to Menu</button>
        </div>
      </div>
    </div>
  `;

  // Add modal to page
  document.body.insertAdjacentHTML("beforeend", modalHTML);

  // Add event listeners
  document.getElementById("viewLeaderboard").addEventListener("click", () => {
    showLeaderboard();
  });

  document.getElementById("playAgain").addEventListener("click", () => {
    window.location.href = "/api/play";
  });

  document.getElementById("backToMenu").addEventListener("click", () => {
    window.location.href = "/";
  });
}

async function showLeaderboard() {
  try {
    const response = await fetch(window.API_ENDPOINTS.LEADERBOARD);
    const result = await response.json();

    if (result.success) {
      displayLeaderboard(result.leaderboard);
    } else {
      console.error("Failed to load leaderboard:", result.error);
    }
  } catch (error) {
    console.error("Error loading leaderboard:", error);
  }
}

function displayLeaderboard(leaderboard) {
  const leaderboardHTML = leaderboard
    .map(
      (entry) => `
    <tr style="border-bottom: 1px solid #34495e;">
      <td style="padding: 0.5rem; text-align: center;">${entry.rank}</td>
      <td style="padding: 0.5rem;">${entry.username}</td>
      <td style="padding: 0.5rem; text-align: center;">${entry.score}</td>
      <td style="padding: 0.5rem; text-align: center;">${entry.completion_time ? formatTime(entry.completion_time) : "--:--"}</td>
    </tr>
  `,
    )
    .join("");

  const modalHTML = `
    <div id="leaderboardModal" style="
      position: fixed;
      top: 0;
      left: 0;
      width: 100%;
      height: 100%;
      background: rgba(0, 0, 0, 0.8);
      display: flex;
      justify-content: center;
      align-items: center;
      z-index: 1001;
    ">
      <div style="
        background: #2c3e50;
        color: white;
        padding: 2rem;
        border-radius: 10px;
        max-width: 600px;
        width: 90%;
        max-height: 80vh;
        overflow-y: auto;
      ">
        <h2 style="color: #ffd700; margin-bottom: 1rem; text-align: center;">🏆 Leaderboard 🏆</h2>
        <table style="width: 100%; border-collapse: collapse;">
          <thead>
            <tr style="background: #34495e;">
              <th style="padding: 0.8rem; text-align: center;">Rank</th>
              <th style="padding: 0.8rem;">Player</th>
              <th style="padding: 0.8rem; text-align: center;">Score</th>
              <th style="padding: 0.8rem; text-align: center;">Time</th>
            </tr>
          </thead>
          <tbody>
            ${leaderboardHTML}
          </tbody>
        </table>
        <div style="text-align: center; margin-top: 1.5rem;">
          <button id="closeLeaderboard" style="
            background: #95a5a6;
            color: white;
            border: none;
            padding: 0.8rem 1.5rem;
            border-radius: 5px;
            cursor: pointer;
            font-size: 1rem;
          ">Close</button>
        </div>
      </div>
    </div>
  `;

  document.body.insertAdjacentHTML("beforeend", modalHTML);

  document.getElementById("closeLeaderboard").addEventListener("click", () => {
    document.getElementById("leaderboardModal").remove();
  });
}
