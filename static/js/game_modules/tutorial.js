let currentTutorialCommand = null;

export function initializeTutorialMode() {
  document.addEventListener("keydown", function (event) {
    if (event.key === window.TUTORIAL_CONFIG.TOGGLE_KEY) {
      toggleTutorialMode();
      event.preventDefault();
    }
  });
}

export function toggleTutorialMode() {
  window.gameState.tutorialMode = !window.gameState.tutorialMode;

  if (window.gameState.tutorialMode) {
    activateTutorialMode();
  } else {
    deactivateTutorialMode();
  }
}

function activateTutorialMode() {
  showTutorialMessage(
    window.TUTORIAL_CONFIG.MESSAGES.ACTIVATED,
    window.TUTORIAL_CONFIG.COLORS.ACTIVATED,
  );
  window.chatModule.addToChatHistory("Tutorial mode activated");

  setTimeout(() => {
    generateRandomTutorialCommand();
  }, window.TUTORIAL_CONFIG.TIMINGS.ACTIVATION_DELAY);
}

function deactivateTutorialMode() {
  currentTutorialCommand = null;
  window.chatModule.addToChatHistory(
    window.TUTORIAL_CONFIG.MESSAGES.DEACTIVATED,
  );
  resetToWelcomeMessage();
}

function resetToWelcomeMessage() {
  const headerInfo = document.querySelector(window.UI_SELECTORS.HEADER_INFO);
  if (!headerInfo) return;

  if (!headerInfo.dataset.originalContent) {
    headerInfo.innerHTML = window.TUTORIAL_CONFIG.MESSAGES.DEFAULT_WELCOME;
  } else {
    headerInfo.innerHTML = headerInfo.dataset.originalContent;
  }
}

export function generateRandomTutorialCommand() {
  if (!window.gameState.tutorialMode) return;

  const randomIndex = Math.floor(
    Math.random() * window.TUTORIAL_COMMANDS.length,
  );
  currentTutorialCommand = window.TUTORIAL_COMMANDS[randomIndex];

  showTutorialMessage(
    currentTutorialCommand.message,
    window.TUTORIAL_CONFIG.COLORS.INSTRUCTION,
  );
}

export function handleTutorialMovement(key) {
  if (!currentTutorialCommand) return;

  if (key === currentTutorialCommand.key) {
    handleCorrectAnswer();
  } else {
    handleWrongAnswer(key);
  }
}

function handleCorrectAnswer() {
  const correctMessage = `✓ CORRECT! ${currentTutorialCommand.message}`;
  showTutorialMessage(correctMessage, window.TUTORIAL_CONFIG.COLORS.CORRECT);
  window.chatModule.addToChatHistory(
    `✓ Correct: ${currentTutorialCommand.message}`,
  );

  setTimeout(() => {
    generateRandomTutorialCommand();
  }, window.TUTORIAL_CONFIG.TIMINGS.NEXT_COMMAND_DELAY);
}

function handleWrongAnswer(pressedKey) {
  const pressedCommand = window.TUTORIAL_COMMANDS.find(
    (cmd) => cmd.key === pressedKey,
  );
  if (!pressedCommand) return;

  const wrongMessage = `✗ Wrong! Expected: ${currentTutorialCommand.message}`;
  showTutorialMessage(wrongMessage, window.TUTORIAL_CONFIG.COLORS.WRONG);

  const chatMessage = `✗ Wrong: Expected ${currentTutorialCommand.message}, pressed ${pressedCommand.message.split(" ")[1]} instead`;
  window.chatModule.addToChatHistory(chatMessage);

  setTimeout(() => {
    showTutorialMessage(
      currentTutorialCommand.message,
      window.TUTORIAL_CONFIG.COLORS.INSTRUCTION,
    );
  }, window.TUTORIAL_CONFIG.TIMINGS.WRONG_ANSWER_DISPLAY);
}

function showTutorialMessage(message, color) {
  const headerInfo = document.querySelector(window.UI_SELECTORS.HEADER_INFO);
  if (!headerInfo) return;

  headerInfo.innerHTML = `<strong style="color: ${color};">${message}</strong>`;
}
