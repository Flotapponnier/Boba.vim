export function initializeMapToggle() {
  document.addEventListener("keydown", function (event) {
    if (event.key === "-") {
      toggleMap();
      event.preventDefault();

      if (!window.gameState.tutorialMode) {
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
          window.chatModule.addToChatHistory(
            `You pressed - to toggle map | ${message}`,
          );

          setTimeout(() => {
            window.feedbackModule.resetHeaderInfo();
          }, 2000);
        }
      } else {
        const mapDisplay = document.getElementById("mapDisplay");
        const isVisible =
          mapDisplay &&
          mapDisplay.style.display !== "none" &&
          mapDisplay.style.display !== "";
        const message = isVisible ? "Map HIDDEN" : "Map SHOWN";
        window.chatModule.addToChatHistory(
          `You pressed - to toggle map | ${message}`,
        );
      }
    }
  });
}

export function toggleMap() {
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
