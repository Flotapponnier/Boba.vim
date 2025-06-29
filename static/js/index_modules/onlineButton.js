export function initializeOnlineButton() {
  const playOnline = document.getElementById("playOnline");

  if (!playOnline) {
    console.log("playOnline button not found");
    return;
  }

  console.log("Initializing online button...");

  playOnline.addEventListener("click", async function () {
    console.log("Online button clicked");

    playOnline.disabled = true;
    playOnline.textContent = "Connecting...";

    try {
      const response = await fetch("/api/playonline", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
      });

      const data = await response.json();
      console.log("Online game response:", data);

      if (data.success) {
        // Handle successful online connection
        alert("Online mode: " + data.message);
      } else {
        alert(data.message || "Online mode not available yet");
      }
    } catch (error) {
      console.error("Error starting online game:", error);
      alert("Failed to connect to online mode. Please try again.");
    } finally {
      playOnline.disabled = false;
      playOnline.textContent = "🧋 Play online";
    }
  });
}
