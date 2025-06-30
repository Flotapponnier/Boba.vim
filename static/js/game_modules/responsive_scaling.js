export function initializeResponsiveScaling() {
  applyScaling();
  
  window.addEventListener('resize', applyScaling);
  
  const resizeObserver = new ResizeObserver(applyScaling);
  const gameBoard = document.querySelector('.game-board');
  if (gameBoard) {
    resizeObserver.observe(gameBoard);
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
  setTimeout(applyScaling, 100);
}