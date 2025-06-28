export function initializeMovement() {
  let lastKeyPressed = null;
  let keyReleased = true;

  document.addEventListener("keydown", function (event) {
    const key = event.key.toLowerCase();

    if (["h", "j", "k", "l"].includes(key)) {
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

    if (["h", "j", "k", "l"].includes(key)) {
      if (lastKeyPressed === key) {
        keyReleased = true;
      }
    }
  });
}

export async function movePlayer(direction) {
  try {
    const response = await fetch("/api/move", {
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
      window.displayModule.updateGameDisplay(result.game_map);
      window.displayModule.updateScore(result.score);

      if (result.pearl_collected && !window.gameState.tutorialMode) {
        console.log("Pearl collected! +100 points");

        const headerInfo = document.querySelector(".header-info");
        if (headerInfo) {
          if (!headerInfo.dataset.originalContent) {
            headerInfo.dataset.originalContent = headerInfo.innerHTML;
          }

          headerInfo.innerHTML = `<strong style="color: #ffd700; animation: pulse 0.5s ease-in-out;">🧋 PEARL COLLECTED! +100 points!</strong>`;
          window.chatModule.addToChatHistory(
            "🧋 PEARL COLLECTED! +100 points!",
          );

          setTimeout(() => {
            const movementMessages = {
              h: "You pressed H to go LEFT ←",
              j: "You pressed J to go DOWN ↓",
              k: "You pressed K to go UP ↑",
              l: "You pressed L to go RIGHT →",
            };
            const message =
              movementMessages[direction] ||
              `You pressed ${direction.toUpperCase()}`;
            headerInfo.innerHTML = `<strong style="color: #ff6b6b;">${message}</strong>`;
          }, 1500);
        }
      } else if (result.pearl_collected && window.gameState.tutorialMode) {
        window.chatModule.addToChatHistory("🧋 PEARL COLLECTED! +100 points!");
      }
    } else {
      console.log("Move blocked:", result.error);

      if (!window.gameState.tutorialMode) {
        const headerInfo = document.querySelector(".header-info");
        if (headerInfo) {
          if (!headerInfo.dataset.originalContent) {
            headerInfo.dataset.originalContent = headerInfo.innerHTML;
          }

          const movementMessages = {
            h: "You pressed H to go LEFT ← - BLOCKED!",
            j: "You pressed J to go DOWN ↓ - BLOCKED!",
            k: "You pressed K to go UP ↑ - BLOCKED!",
            l: "You pressed L to go RIGHT → - BLOCKED!",
          };

          const message =
            movementMessages[direction] ||
            `You pressed ${direction.toUpperCase()} - BLOCKED!`;
          headerInfo.innerHTML = `<strong style="color: #e74c3c;">${message}</strong>`;
          window.chatModule.addToChatHistory(message);
        }
      } else {
        const movementMessages = {
          h: "You pressed H to go LEFT ← - BLOCKED!",
          j: "You pressed J to go DOWN ↓ - BLOCKED!",
          k: "You pressed K to go UP ↑ - BLOCKED!",
          l: "You pressed L to go RIGHT → - BLOCKED!",
        };
        const message =
          movementMessages[direction] ||
          `You pressed ${direction.toUpperCase()} - BLOCKED!`;
        window.chatModule.addToChatHistory(message);
      }
    }
  } catch (error) {
    console.error("Error moving player:", error);
  }
}
