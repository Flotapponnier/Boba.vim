export function initializePlayButton() {
  const playButton = document.getElementById("playButton");

  if (!playButton) {
    console.log("playButton not found");
    return;
  }

  console.log("Initializing play button...");

  playButton.addEventListener("click", async function () {
    console.log("Play button clicked");

    const username = await promptForUsername();
    if (!username) return;

    playButton.disabled = true;
    playButton.textContent = "🚀 Starting...";

    try {
      // Set username before starting game
      await setUsername(username);
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

async function promptForUsername() {
  const username = prompt("Enter your username for the leaderboard:");
  if (!username || username.trim().length < 2) {
    alert("Username must be at least 2 characters long.");
    return null;
  }
  if (username.length > 50) {
    alert("Username must be less than 50 characters.");
    return null;
  }
  return username.trim();
}

async function setUsername(username) {
  try {
    const response = await fetch("/api/set-username", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ username: username }),
    });

    const result = await response.json();
    if (!result.success) {
      throw new Error(result.error);
    }
    return result;
  } catch (error) {
    console.error("Error setting username:", error);
    throw error;
  }
}
