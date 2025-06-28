export function initializeMapToggle() {
  document.addEventListener("keydown", function (event) {
    if (event.key === window.MAP_CONFIG.TOGGLE_KEY) {
      handleMapToggle();
      event.preventDefault();
    }
  });
}

function handleMapToggle() {
  const wasVisible = isMapVisible();
  toggleMap();
  const isNowVisible = isMapVisible();

  const statusMessage = isNowVisible
    ? window.MAP_CONFIG.MESSAGES.SHOWN
    : window.MAP_CONFIG.MESSAGES.HIDDEN;
  const fullMessage = `${window.MAP_CONFIG.MESSAGES.TOGGLE} | ${statusMessage}`;

  window.chatModule.addToChatHistory(fullMessage);

  if (!window.gameState.tutorialMode) {
    showMapToggleFeedback(fullMessage);
  }
}

function isMapVisible() {
  const mapDisplay = document.getElementById(
    window.MAP_CONFIG.DISPLAY_ELEMENT_ID,
  );
  return (
    mapDisplay &&
    mapDisplay.style.display !== "none" &&
    mapDisplay.style.display !== ""
  );
}

function showMapToggleFeedback(message) {
  const headerInfo = document.querySelector(window.UI_SELECTORS.HEADER_INFO);
  if (!headerInfo) return;

  if (!headerInfo.dataset.originalContent) {
    headerInfo.dataset.originalContent = headerInfo.innerHTML;
  }

  headerInfo.innerHTML = `<strong style="color: ${window.MAP_CONFIG.FEEDBACK_COLOR};">${message}</strong>`;

  setTimeout(() => {
    window.feedbackModule.resetHeaderInfo();
  }, window.MAP_CONFIG.FEEDBACK_TIMEOUT);
}

export function toggleMap() {
  const mapDisplay = document.getElementById(
    window.MAP_CONFIG.DISPLAY_ELEMENT_ID,
  );
  if (!mapDisplay) return;

  const isCurrentlyVisible =
    mapDisplay.style.display !== "none" && mapDisplay.style.display !== "";

  mapDisplay.style.display = isCurrentlyVisible ? "none" : "block";
}
