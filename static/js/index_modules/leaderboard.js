export function initializeLeaderboardButton() {
  const leaderboardButton = document.getElementById("leaderboardButton");

  if (!leaderboardButton) {
    console.log("leaderboardButton not found");
    return;
  }

  console.log("Initializing leaderboard button...");

  leaderboardButton.addEventListener("mouseenter", function () {
    leaderboardButton.textContent = "🧋";
  });

  leaderboardButton.addEventListener("mouseleave", function () {
    leaderboardButton.textContent = "🧋Leaderboard";
  });

  leaderboardButton.addEventListener("click", async function () {
    console.log("Leaderboard button clicked");

    leaderboardButton.disabled = true;
    leaderboardButton.textContent = "Loading...";

    try {
      const response = await fetch("/api/leaderboard");
      const result = await response.json();

      console.log("Leaderboard response:", result);

      if (result.success) {
        showLeaderboardModal(result.leaderboard);
      } else {
        alert("Failed to load leaderboard: " + result.error);
      }
    } catch (error) {
      console.error("Error loading leaderboard:", error);
      alert("Failed to load leaderboard. Please try again.");
    } finally {
      leaderboardButton.disabled = false;
      leaderboardButton.textContent = "🏆 Leaderboard";
    }
  });
}

function showLeaderboardModal(leaderboard) {
  console.log("Showing leaderboard modal with data:", leaderboard);

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
        box-shadow: 0 10px 30px rgba(0,0,0,0.5);
      ">
        <h2 style="color: #ffd700; margin-bottom: 1rem; text-align: center;">🏆 Leaderboard 🏆</h2>
        <table style="width: 100%; border-collapse: collapse;">
          <thead>
            <tr style="background: #34495e;">
              <th style="padding: 0.8rem; text-align: center;">Rank</th>
              <th style="padding: 0.8rem;">Player</th>
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
            transition: background 0.3s ease;
          " onmouseover="this.style.background='#7f8c8d'" onmouseout="this.style.background='#95a5a6'">Close</button>
        </div>
      </div>
    </div>
  `;

  document.body.insertAdjacentHTML("beforeend", modalHTML);

  // Close modal events
  const modal = document.getElementById("leaderboardModal");
  document.getElementById("closeLeaderboard").addEventListener("click", () => {
    modal.remove();
  });

  // Close on outside click
  modal.addEventListener("click", (e) => {
    if (e.target === modal) {
      modal.remove();
    }
  });

  // Close on Escape key
  document.addEventListener("keydown", function escapeHandler(e) {
    if (e.key === "Escape") {
      modal.remove();
      document.removeEventListener("keydown", escapeHandler);
    }
  });
}

function formatTime(seconds) {
  const minutes = Math.floor(seconds / 60);
  const remainingSeconds = seconds % 60;
  return `${minutes}:${remainingSeconds.toString().padStart(2, "0")}`;
}
