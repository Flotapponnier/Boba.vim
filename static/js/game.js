document.addEventListener("DOMContentLoaded", function () {
  initializeGame();
  initializeBackToMenuButton();
  initializeMovement();
  initializeMapToggle();
  initializeTutorialMode();
});

function initializeGame() {
  console.log("🎮 Boba.vim Text Adventure loaded - Python backend");
}

let tutorialMode = false;
let currentTutorialCommand = null;
let tutorialCommands = [
  { key: "h", message: "Press H to go LEFT ←" },
  { key: "j", message: "Press J to go DOWN ↓" },
  { key: "k", message: "Press K to go UP ↑" },
  { key: "l", message: "Press L to go RIGHT →" },
];

function initializeTutorialMode() {
  document.addEventListener("keydown", function (event) {
    if (event.key === "+") {
      toggleTutorialMode();
      event.preventDefault();
    }
  });
}

function toggleTutorialMode() {
  tutorialMode = !tutorialMode;
  const headerInfo = document.querySelector(".header-info");

  if (tutorialMode) {
    headerInfo.innerHTML = `<strong style="color: #9b59b6;">TUTORIAL MODE ACTIVATED</strong>`;
    setTimeout(() => {
      generateRandomTutorialCommand();
    }, 1000);
  } else {
    currentTutorialCommand = null;
    if (!headerInfo.dataset.originalContent) {
      headerInfo.innerHTML =
        "Welcome! Use HJKL to move | Press - to show map | Press + to activate tutorial";
    } else {
      headerInfo.innerHTML = headerInfo.dataset.originalContent;
    }
  }
}

function generateRandomTutorialCommand() {
  if (!tutorialMode) return;

  const randomIndex = Math.floor(Math.random() * tutorialCommands.length);
  currentTutorialCommand = tutorialCommands[randomIndex];

  const headerInfo = document.querySelector(".header-info");
  headerInfo.innerHTML = `<strong style="color: #3498db;">${currentTutorialCommand.message}</strong>`;
}

function initializeBackToMenuButton() {
  const menuButton = document.getElementById("backMenu");
  if (!menuButton) return;
  menuButton.addEventListener("click", function () {
    menuButton.disabled = true;
    menuButton.textContent = "Going back to menu ...";
    try {
      window.location.href = "/";
    } catch (error) {
      console.error("Error going back to menu:", error);
      alert("Error going back to menu");
      menuButton.disabled = false;
      menuButton.textContent = "← Menu";
    }
  });
}

function initializeMovement() {
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

      if (tutorialMode) {
        handleTutorialMovement(key);
      } else {
        showMovementFeedback(key);
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

function handleTutorialMovement(key) {
  const headerInfo = document.querySelector(".header-info");

  if (currentTutorialCommand && key === currentTutorialCommand.key) {
    headerInfo.innerHTML = `<strong style="color: #27ae60;">✓ CORRECT! ${currentTutorialCommand.message}</strong>`;

    setTimeout(() => {
      generateRandomTutorialCommand();
    }, 1000);
  } else if (currentTutorialCommand) {
    const correctCommand = tutorialCommands.find((cmd) => cmd.key === key);
    if (correctCommand) {
      headerInfo.innerHTML = `<strong style="color: #e74c3c;">✗ Wrong! Expected: ${currentTutorialCommand.message}</strong>`;

      setTimeout(() => {
        headerInfo.innerHTML = `<strong style="color: #3498db;">${currentTutorialCommand.message}</strong>`;
      }, 1500);
    }
  }
}

function showMovementFeedback(key) {
  const movementMessages = {
    h: "You pressed H to go LEFT ←",
    j: "You pressed J to go DOWN ↓",
    k: "You pressed K to go UP ↑",
    l: "You pressed L to go RIGHT →",
  };

  const headerInfo = document.querySelector(".header-info");
  if (headerInfo) {
    if (!headerInfo.dataset.originalContent) {
      headerInfo.dataset.originalContent = headerInfo.innerHTML;
    }

    const message = movementMessages[key] || `You pressed ${key.toUpperCase()}`;
    headerInfo.innerHTML = `<strong style="color: #ff6b6b;">${message}</strong>`;

    headerInfo.style.animation = "pulse 0.3s ease-in-out";

    setTimeout(() => {
      headerInfo.style.animation = "";
    }, 300);
  }
}

function resetHeaderInfo() {
  const headerInfo = document.querySelector(".header-info");
  if (headerInfo && headerInfo.dataset.originalContent) {
    headerInfo.innerHTML = headerInfo.dataset.originalContent;
  }
}

function initializeMapToggle() {
  document.addEventListener("keydown", function (event) {
    if (event.key === "-") {
      toggleMap();
      event.preventDefault();

      if (!tutorialMode) {
        const headerInfo = document.querySelector(".header-info");
        if (headerInfo) {
          if (!headerInfo.dataset.originalContent) {
            headerInfo.dataset.originalContent = headerInfo.innerHTML;
          }

          const mapDisplay = document.getElementById("mapDisplay");
          const isVisible =
            mapDisplay &&
            mapDisplay.style.display !== "none" &&
            mapDisplay.style.display !== "";
          const message = isVisible ? "Map HIDDEN" : "Map SHOWN";

          headerInfo.innerHTML = `<strong style="color: #4ecdc4;">You pressed - to toggle map | ${message}</strong>`;

          setTimeout(() => {
            resetHeaderInfo();
          }, 2000);
        }
      }
    }
  });
}

function toggleMap() {
  const mapDisplay = document.getElementById("mapDisplay");
  if (mapDisplay) {
    if (
      mapDisplay.style.display === "none" ||
      mapDisplay.style.display === ""
    ) {
      mapDisplay.style.display = "block";
    } else {
      mapDisplay.style.display = "none";
    }
  }
}

async function movePlayer(direction) {
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
      updateGameDisplay(result.game_map);
      updateScore(result.score);

      if (result.pearl_collected && !tutorialMode) {
        console.log("Pearl collected! +100 points");

        const headerInfo = document.querySelector(".header-info");
        if (headerInfo) {
          if (!headerInfo.dataset.originalContent) {
            headerInfo.dataset.originalContent = headerInfo.innerHTML;
          }

          headerInfo.innerHTML = `<strong style="color: #ffd700; animation: pulse 0.5s ease-in-out;">🧋 PEARL COLLECTED! +100 points!</strong>`;

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
      }
    } else {
      console.log("Move blocked:", result.error);

      if (!tutorialMode) {
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
        }
      }
    }
  } catch (error) {
    console.error("Error moving player:", error);
  }
}

function updateGameDisplay(gameMap) {
  const keys = document.querySelectorAll(".key");

  keys.forEach((key) => {
    const row = parseInt(key.getAttribute("data-row"));
    const col = parseInt(key.getAttribute("data-col"));
    const oldMapValue = parseInt(key.getAttribute("data-map"));
    const newMapValue = gameMap[row][col];

    if (oldMapValue !== newMapValue) {
      key.setAttribute("data-map", newMapValue);

      const existingBoba = key.querySelector(".boba-character");
      const existingPearl = key.querySelector(".pearl");
      if (existingBoba) existingBoba.remove();
      if (existingPearl) existingPearl.remove();

      const keyTop = key.querySelector(".key-top");

      if (newMapValue === 1) {
        const bobaDiv = document.createElement("div");
        bobaDiv.className = "boba-character";
        bobaDiv.innerHTML = `
          <div class="boba-shadow"></div>
          <img src="/static/sprites/boba.png" alt="Boba" class="boba-sprite">
        `;
        keyTop.appendChild(bobaDiv);
      } else if (newMapValue === 3) {
        const pearlDiv = document.createElement("div");
        pearlDiv.className = "pearl";
        pearlDiv.innerHTML = `
          <div class="pearl-shadow"></div>
          <img src="/static/sprites/pearl.png" alt="Pearl" class="pearl-sprite">
        `;
        keyTop.appendChild(pearlDiv);
      }
    }
  });

  updateDebugDisplay(gameMap);
}

function updateDebugDisplay(gameMap) {
  const mapGrid = document.getElementById("mapGrid");
  if (mapGrid) {
    let html = "";
    for (let row = 0; row < gameMap.length; row++) {
      html += '<div class="map-row">';
      for (let col = 0; col < gameMap[row].length; col++) {
        const value = gameMap[row][col];
        html += `<span class="map-cell map-${value}">${value}</span>`;
      }
      html += "</div>";
    }
    mapGrid.innerHTML = html;
  }
}

function updateScore(score) {
  const scoreElement = document.getElementById("score");
  if (scoreElement) {
    scoreElement.textContent = score;

    scoreElement.style.color = "#4ecdc4";
    scoreElement.style.transform = "scale(1.2)";

    setTimeout(() => {
      scoreElement.style.color = "#333";
      scoreElement.style.transform = "scale(1)";
    }, 300);
  }
}
