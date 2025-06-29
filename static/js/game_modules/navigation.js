export function initializeBackToMenuButton() {
  const menuButton = document.getElementById(
    window.UI_SELECTORS.BACK_MENU_BUTTON.replace("#", ""),
  );
  if (!menuButton) return;

  menuButton.addEventListener("click", function () {
    menuButton.disabled = true;
    try {
      window.location.href = "/";
    } catch (error) {
      console.error("Error going back to menu:", error);
      alert("Error going back to menu");
      menuButton.disabled = false;
      menuButton.textContent = "← Menu";
    }
  });

  menuButton.addEventListener("mouseenter", function () {
    menuButton.textContent = "🧋";
  });

  menuButton.addEventListener("mouseleave", function () {
    menuButton.textContent = "Menu";
  });
}
