document.addEventListener("DOMContentLoaded", function () {
  initializePlayButton();
  initializeTutorialButton();
  initializeOnlineButton();
});

function initializePlayButton() {
  const playButton = document.getElementById("playButton");

  if (!playButton) return;

  playButton.addEventListener("click", async function () {
    playButton.disabled = true;
    playButton.textContent = "🚀 Starting...";

    try {
      window.location.href = "/api/play";
    } catch (error) {
      console.error("Error starting game:", error);
      alert("Failed to start game. Please try again.");
    } finally {
      playButton.disabled = false;
      playButton.textContent = "🧋 Play";
    }
  });
}

function initializeTutorialButton() {
  const playTutorialButton = document.getElementById("tutorialplayButton");

  if (!playTutorialButton) return;

  playTutorialButton.addEventListener("click", async function () {
    playTutorialButton.disabled = true;
    playTutorialButton.textContent = "To be implemented..";

    try {
      const response = await fetch("/api/playtutorial");
      const data = await response.json();

      console.log("Game response: ", data);
      alert(data.message);
    } catch (error) {
      console.error("Error starting game:", error);
      alert("Failed to start game. Please try again.");
    } finally {
      playTutorialButton.disabled = false;
      playTutorialButton.textContent = "🧋 Play with tutorial";
    }
  });
}

function initializeOnlineButton() {
  const playOnline = document.getElementById("playOnline");

  if (!playOnline) return;

  playOnline.addEventListener("click", async function () {
    playOnline.disabled = true;
    playOnline.textContent = "To be implemented..";

    try {
      const response = await fetch("/api/playonline");
      const data = await response.json();
      console.log("Game response:", data);
      alert(data.message);
    } catch (error) {
      console.error("Error starting game", error);
      alert("Failed to start game. Please try again");
    } finally {
      playOnline.disabled = false;
      playOnline.textContent = "🧋 Play online";
    }
  });
}
