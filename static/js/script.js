//Playbutton
document.addEventListener("DOMContentLoaded", function () {
  const playButton = document.getElementById("playButton");

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
});

//TutorialButton
document.addEventListener("DOMContentLoaded", function () {
  const playTutorialButton = document.getElementById("tutorialplayButton");

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
});

//Online button
document.addEventListener("DOMContentLoaded", function () {
  const playOnline = document.getElementById("playOnline");
  playOnline.addEventListener("click", async function () {
    playOnline.disabled = true;
    playOnline.textcontent = "To be implemented..";
    try {
      const response = await fetch("/api/playonline");
      const data = await response.json();
      console.log("Game response:", data);
      alert(data.message);
    } catch (error) {
      console.error("Error starting game", error);
      alert("Failed to start game. Please try again");
    } finally {
      playOnline.disabled = False;
      playOnline.textContent = "🧋 Play online";
    }
  });
});

document.addEventListener("DOMContentLoaded", function () {
  const menuButton = document.getElementById("backMenu");

  menuButton.addEventListener("click", function () {
    menuButton.disabled = true;
    menuButton.textContent = "Going back to menu ...";
    try {
      window.location.href = "/";
    } catch (error) {
      console.error("Error going back to menu:", error);
      alert("Error going back to menu");
      menuButton.disabled = false;
      menuButton.textContent = "Menu";
    }
  });
});
