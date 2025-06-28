function initializeBackToMenuButton() {
  const menuButton = document.getElementById("backMenu");

  if (!menuButton) return;

  menuButton.addEventListener("click", function () {
    menuButton.disabled = true;
    menuButton.textContent = "Going back to menu ...";

    try {
      window.location.href = "/";
    } catch (error) {
      console.error("Error going back to menu:", error);
      alert("Error going back to menu");
      menuButton.disabled = false;
      menuButton.textContent = "← Menu";
    }
  });
}
