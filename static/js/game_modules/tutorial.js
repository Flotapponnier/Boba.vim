const tutorialCommands = [
  { key: "h", message: "Press H to go LEFT ←" },
  { key: "j", message: "Press J to go DOWN ↓" },
  { key: "k", message: "Press K to go UP ↑" },
  { key: "l", message: "Press L to go RIGHT →" },
];

let currentTutorialCommand = null;

export function initializeTutorialMode() {
  document.addEventListener("keydown", function (event) {
    if (event.key === "+") {
      toggleTutorialMode();
      event.preventDefault();
    }
  });
}

export function toggleTutorialMode() {
  window.gameState.tutorialMode = !window.gameState.tutorialMode;
  const headerInfo = document.querySelector(".header-info");

  if (window.gameState.tutorialMode) {
    headerInfo.innerHTML = `<strong style="color: #9b59b6;">TUTORIAL MODE ACTIVATED</strong>`;
    window.chatModule.addToChatHistory("Tutorial mode activated");
    setTimeout(() => {
      generateRandomTutorialCommand();
    }, 1000);
  } else {
    currentTutorialCommand = null;
    window.chatModule.addToChatHistory("Tutorial mode deactivated");
    if (!headerInfo.dataset.originalContent) {
      headerInfo.innerHTML =
        "Welcome! Use HJKL to move | Press - to show map | Press + to activate tutorial";
    } else {
      headerInfo.innerHTML = headerInfo.dataset.originalContent;
    }
  }
}

export function generateRandomTutorialCommand() {
  if (!window.gameState.tutorialMode) return;

  const randomIndex = Math.floor(Math.random() * tutorialCommands.length);
  currentTutorialCommand = tutorialCommands[randomIndex];

  const headerInfo = document.querySelector(".header-info");
  headerInfo.innerHTML = `<strong style="color: #3498db;">${currentTutorialCommand.message}</strong>`;
}

export function handleTutorialMovement(key) {
  const headerInfo = document.querySelector(".header-info");

  if (currentTutorialCommand && key === currentTutorialCommand.key) {
    headerInfo.innerHTML = `<strong style="color: #27ae60;">✓ CORRECT! ${currentTutorialCommand.message}</strong>`;
    window.chatModule.addToChatHistory(
      `✓ Correct: ${currentTutorialCommand.message}`,
    );

    setTimeout(() => {
      generateRandomTutorialCommand();
    }, 1000);
  } else if (currentTutorialCommand) {
    const correctCommand = tutorialCommands.find((cmd) => cmd.key === key);
    if (correctCommand) {
      headerInfo.innerHTML = `<strong style="color: #e74c3c;">✗ Wrong! Expected: ${currentTutorialCommand.message}</strong>`;
      window.chatModule.addToChatHistory(
        `✗ Wrong: Expected ${currentTutorialCommand.message}, pressed ${correctCommand.message.split(" ")[1]} instead`,
      );

      setTimeout(() => {
        headerInfo.innerHTML = `<strong style="color: #3498db;">${currentTutorialCommand.message}</strong>`;
      }, 1500);
    }
  }
}
