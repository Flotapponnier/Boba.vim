export function initializeMovement() {
  let lastKeyPressed = null;
  let keyReleased = true;

  document.addEventListener("keydown", function (event) {
    const key = event.key.toLowerCase();

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
    const key = event.key.toLowerCase();

    if (window.VALID_MOVEMENT_KEYS.includes(key)) {
      if (lastKeyPressed === key) {
        keyReleased = true;
      }
    }
  });
}

export async function movePlayer(direction) {
  try {
    const response = await fetch(window.API_ENDPOINTS.MOVE, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        direction: window.MOVEMENT_KEYS[direction].direction,
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
  }
}

function handleSuccessfulMove(result, direction) {
  window.displayModule.updateGameDisplay(result.game_map);
  window.displayModule.updateScore(result.score);

  if (result.pearl_collected) {
    handlePearlCollection(direction);
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
