let currentScore = 0;

document.addEventListener("DOMContentLoaded", function () {
  initializeGame();
  initializeBackToMenuButton();
});

function initializeGame() {
  console.log("Game initialized");
  updateScore(0);
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

function updateScore(newScore) {
  currentScore = newScore;
  const scoreElement = document.getElementById("score");
  if (scoreElement) {
    scoreElement.textContent = currentScore;
  }
}

function addScore(points) {
  updateScore(currentScore + points);
}

function getScore() {
  return currentScore;
}

function resetScore() {
  updateScore(0);
}

if (typeof window !== "undefined") {
  window.gameScore = {
    update: updateScore,
    add: addScore,
    get: getScore,
    reset: resetScore,
  };
}
