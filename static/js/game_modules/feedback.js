export function showMovementFeedback(key) {
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

    window.chatModule.addToChatHistory(message);

    headerInfo.style.animation = "pulse 0.3s ease-in-out";

    setTimeout(() => {
      headerInfo.style.animation = "";
    }, 300);
  }
}

export function resetHeaderInfo() {
  const headerInfo = document.querySelector(".header-info");
  if (headerInfo && headerInfo.dataset.originalContent) {
    headerInfo.innerHTML = headerInfo.dataset.originalContent;
  }
}
