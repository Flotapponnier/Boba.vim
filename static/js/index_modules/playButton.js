import { getSelectedCharacter } from './characterSelection.js';

export function initializePlayButton() {
  const playButton = document.getElementById("playButton");

  if (!playButton) {
    console.log("playButton not found");
    return;
  }

  console.log("Initializing play button...");

  playButton.addEventListener("mouseenter", function () {
    playButton.textContent = "🧋";
  });

  playButton.addEventListener("mouseleave", function () {
    playButton.textContent = "🧋 Play";
  });

  playButton.addEventListener("click", function () {
    console.log("Play button clicked");

    playButton.disabled = true;
    playButton.textContent = "🚀 Starting...";

    try {
      // Get selected character and include in URL
      const selectedCharacter = getSelectedCharacter();
      window.location.href = `/api/play?character=${encodeURIComponent(selectedCharacter)}`;
    } catch (error) {
      console.error("Error starting game:", error);
      alert("Failed to start game. Please try again.");
      playButton.disabled = false;
      playButton.textContent = "🧋 Play";
    }
  });
}

