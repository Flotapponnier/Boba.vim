let selectedCharacter = 'boba';

export function initializeCharacterSelection() {
  console.log("Initializing character selection...");
  
  const characterBoxes = document.querySelectorAll('.character-box');
  const modal = document.getElementById('registrationModal');
  
  characterBoxes.forEach(box => {
    box.addEventListener('click', () => {
      const character = box.dataset.character;
      
      if (box.classList.contains('locked')) {
        showRegistrationModal();
        return;
      }
      
      // Allow selection of any unlocked character (including currently selected)
      if (box.classList.contains('unlocked') || box.classList.contains('selected')) {
        selectCharacter(character);
      }
    });
  });
  
  window.closeRegistrationModal = () => {
    modal.classList.add('hidden');
  };
  
  modal.addEventListener('click', (e) => {
    if (e.target === modal) {
      closeRegistrationModal();
    }
  });
}

function selectCharacter(character) {
  document.querySelectorAll('.character-box').forEach(box => {
    box.classList.remove('selected');
  });
  
  const selectedBox = document.querySelector(`[data-character="${character}"]`);
  if (selectedBox) {
    selectedBox.classList.add('selected');
    selectedCharacter = character;
    console.log(`Character selected: ${character}`);
  }
}

function showRegistrationModal() {
  const modal = document.getElementById('registrationModal');
  modal.classList.remove('hidden');
}

export function getSelectedCharacter() {
  return selectedCharacter;
}