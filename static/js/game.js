document.addEventListener("DOMContentLoaded", function () {
  initializeGame();
  initializeBackToMenuButton();
  initializeMovement();
  initializeMapToggle();
});

function initializeGame() {
  console.log("🎮 Boba.vim Text Adventure loaded - Python backend");
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
  document.addEventListener("keydown", function (event) {
    const key = event.key.toLowerCase();

    if (["h", "j", "k", "l"].includes(key)) {
      movePlayer(key);
      event.preventDefault();
    }
  });
}

function initializeMapToggle() {
  document.addEventListener("keydown", function (event) {
    if (event.key === "-") {
      toggleMap();
      event.preventDefault();
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

      if (result.pearl_collected) {
        console.log("Pearl collected! +100 points");
      }
    } else {
      console.log("Move blocked:", result.error);
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
