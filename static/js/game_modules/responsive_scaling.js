let scalingEnabled = true;
let hasInitialScale = false;

export function initializeResponsiveScaling() {
  applyInitialScaling();
  
  window.addEventListener('resize', debounce(handleResize, 150));
}

function handleResize() {
  if (scalingEnabled) {
    applyScaling();
  }
}

function applyInitialScaling() {
  if (!hasInitialScale) {
    applyScaling();
    hasInitialScale = true;
  }
}

function applyScaling() {
  const gameBoard = document.querySelector('.game-board');
  const keyboardMap = document.querySelector('.keyboard-map');
  
  if (!gameBoard || !keyboardMap) return;
  
  keyboardMap.style.transform = 'scale(1)';
  
  requestAnimationFrame(() => {
    const boardRect = gameBoard.getBoundingClientRect();
    const mapRect = keyboardMap.getBoundingClientRect();
    
    const padding = 80;
    const availableWidth = boardRect.width - padding;
    const availableHeight = boardRect.height - padding;
    
    if (availableWidth <= 0 || availableHeight <= 0) return;
    
    const scaleX = availableWidth / mapRect.width;
    const scaleY = availableHeight / mapRect.height;
    
    const minScale = 0.3;
    const maxScale = 1.2;
    const scale = Math.max(minScale, Math.min(scaleX, scaleY, maxScale));
    
    if (scale !== 1) {
      keyboardMap.style.transform = `scale(${scale})`;
    }
  });
}

export function updateScalingAfterMapChange() {
  // Only apply scaling if we haven't scaled yet (initial game load)
  if (!hasInitialScale) {
    setTimeout(applyInitialScaling, 100);
  }
}

export function disableScaling() {
  scalingEnabled = false;
  // Also disable the scaling transition to prevent any visual movement
  const keyboardMap = document.querySelector('.keyboard-map');
  if (keyboardMap) {
    keyboardMap.style.transition = 'none';
  }
}

export function enableScaling() {
  scalingEnabled = true;
  // Re-enable transition
  const keyboardMap = document.querySelector('.keyboard-map');
  if (keyboardMap) {
    keyboardMap.style.transition = 'transform 0.3s ease';
  }
  if (scalingEnabled) {
    applyScaling();
  }
}

function debounce(func, wait) {
  let timeout;
  return function executedFunction(...args) {
    const later = () => {
      clearTimeout(timeout);
      func(...args);
    };
    clearTimeout(timeout);
    timeout = setTimeout(later, wait);
  };
}