document.addEventListener("DOMContentLoaded", function () {
  initializePlayButton();
  initializeTutorialButton();
  initializeOnlineButton();
  initializeLeaderboardButton();
  initializeUsernameInput();
});

function initializePlayButton() {
  const playButton = document.getElementById("playButton");

  if (!playButton) return;

  playButton.addEventListener("click", async function () {
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

function initializeLeaderboardButton() {
  const leaderboardButton = document.getElementById("leaderboardButton");

  if (!leaderboardButton) return;

  leaderboardButton.addEventListener("click", async function () {
    try {
      const response = await fetch("/api/leaderboard");
      const result = await response.json();

      if (result.success) {
        showLeaderboardModal(result.leaderboard);
      } else {
        alert("Failed to load leaderboard: " + result.error);
      }
    } catch (error) {
      console.error("Error loading leaderboard:", error);
      alert("Failed to load leaderboard. Please try again.");
    }
  });
}

function showLeaderboardModal(leaderboard) {
  if (leaderboard.length === 0) {
    alert("No scores yet! Be the first to complete the game!");
    return;
  }

  const leaderboardHTML = leaderboard
    .map(
      (entry) => `
    <tr style="border-bottom: 1px solid #34495e;">
      <td style="padding: 0.5rem; text-align: center; font-weight: bold;">${entry.rank}</td>
      <td style="padding: 0.5rem;">${entry.username}</td>
      <td style="padding: 0.5rem; text-align: center;">${entry.score}</td>
      <td style="padding: 0.5rem; text-align: center;">${entry.completion_time ? formatTime(entry.completion_time) : "--:--"}</td>
    </tr>
  `,
    )
    .join("");

  const modalHTML = `
    <div id="leaderboardModal" style="
      position: fixed;
      top: 0;
      left: 0;
      width: 100%;
      height: 100%;
      background: rgba(0, 0, 0, 0.8);
      display: flex;
      justify-content: center;
      align-items: center;
      z-index: 1000;
    ">
      <div style="
        background: #2c3e50;
        color: white;
        padding: 2rem;
        border-radius: 10px;
        max-width: 600px;
        width: 90%;
        max-height: 80vh;
        overflow-y: auto;
      ">
        <h2 style="color: #ffd700; margin-bottom: 1rem; text-align: center;">🏆 Leaderboard 🏆</h2>
        <table style="width: 100%; border-collapse: collapse;">
          <thead>
            <tr style="background: #34495e;">
              <th style="padding: 0.8rem; text-align: center;">Rank</th>
              <th style="padding: 0.8rem;">Player</th>
              <th style="padding: 0.8rem; text-align: center;">Score</th>
              <th style="padding: 0.8rem; text-align: center;">Time</th>
            </tr>
          </thead>
          <tbody>
            ${leaderboardHTML}
          </tbody>
        </table>
        <div style="text-align: center; margin-top: 1.5rem;">
          <button id="closeLeaderboard" style="
            background: #95a5a6;
            color: white;
            border: none;
            padding: 0.8rem 1.5rem;
            border-radius: 5px;
            cursor: pointer;
            font-size: 1rem;
          ">Close</button>
        </div>
      </div>
    </div>
  `;

  document.body.insertAdjacentHTML("beforeend", modalHTML);

  document.getElementById("closeLeaderboard").addEventListener("click", () => {
    document.getElementById("leaderboardModal").remove();
  });
}

function formatTime(seconds) {
  const minutes = Math.floor(seconds / 60);
  const remainingSeconds = seconds % 60;
  return `${minutes}:${remainingSeconds.toString().padStart(2, "0")}`;
}

function initializeUsernameInput() {
  // Add any username input initialization logic here if needed
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
