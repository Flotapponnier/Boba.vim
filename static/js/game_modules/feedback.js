export function showMovementFeedback(key) {
  const headerInfo = document.querySelector(window.UI_SELECTORS.HEADER_INFO);
  if (!headerInfo) return;

  storeOriginalHeaderContent(headerInfo);

  const message = window.MOVEMENT_MESSAGES[key] || createFallbackMessage(key);

  updateHeaderWithMessage(headerInfo, message);
  window.chatModule.addToChatHistory(message);

  addPulseAnimation(headerInfo);
}

export function resetHeaderInfo() {
  const headerInfo = document.querySelector(window.UI_SELECTORS.HEADER_INFO);
  if (headerInfo && headerInfo.dataset.originalContent) {
    headerInfo.innerHTML = headerInfo.dataset.originalContent;
  }
}

function storeOriginalHeaderContent(headerInfo) {
  if (!headerInfo.dataset.originalContent) {
    headerInfo.dataset.originalContent = headerInfo.innerHTML;
  }
}

function createFallbackMessage(key) {
  // Handle character search motions with embedded character
  if (key.startsWith('find_char_forward_')) {
    const char = key.slice(18);
    return `FIND CHAR '${char}' →`;
  }
  if (key.startsWith('find_char_backward_')) {
    const char = key.slice(19);
    return `FIND CHAR '${char}' ←`;
  }
  if (key.startsWith('till_char_forward_')) {
    const char = key.slice(18);
    return `TILL CHAR '${char}' →`;
  }
  if (key.startsWith('till_char_backward_')) {
    const char = key.slice(19);
    return `TILL CHAR '${char}' ←`;
  }
  
  return `You pressed ${key.toUpperCase()}`;
}

function updateHeaderWithMessage(headerInfo, message) {
  headerInfo.innerHTML = `<strong style="color: ${window.FEEDBACK_CONFIG.MOVEMENT_COLOR};">${message}</strong>`;
}

function addPulseAnimation(headerInfo) {
  headerInfo.style.animation = window.FEEDBACK_CONFIG.ANIMATION_TYPE;

  setTimeout(() => {
    headerInfo.style.animation = "";
  }, window.FEEDBACK_CONFIG.ANIMATION_DURATION);
}
