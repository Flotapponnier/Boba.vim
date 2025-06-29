let selectedCharacter = 'boba';
let currentUser = null;

export function initializeCharacterSelection() {
  console.log("Initializing character selection...");
  
  const characterBoxes = document.querySelectorAll('.character-box');
  const modal = document.getElementById('authModal');
  const registrationForm = document.getElementById('registrationForm');
  const loginForm = document.getElementById('loginForm');
  const logoutButton = document.getElementById('logoutButton');
  const loginButton = document.getElementById('loginButton');
  const registerButton = document.getElementById('registerButton');
  const switchToLogin = document.getElementById('switchToLogin');
  const switchToRegister = document.getElementById('switchToRegister');
  
  // Check current user status
  checkAuthStatus();
  
  characterBoxes.forEach(box => {
    box.addEventListener('click', () => {
      const character = box.dataset.character;
      
      if (box.classList.contains('locked')) {
        showAuthModal('register');
        return;
      }
      
      // Allow selection of any unlocked character (including currently selected)
      if (box.classList.contains('unlocked') || box.classList.contains('selected')) {
        selectCharacter(character);
      }
    });
  });
  
  // Handle form submissions
  if (registrationForm) {
    registrationForm.addEventListener('submit', handleRegistration);
  }
  
  if (loginForm) {
    loginForm.addEventListener('submit', handleLogin);
  }
  
  // Handle logout
  if (logoutButton) {
    logoutButton.addEventListener('click', handleLogout);
  }
  
  // Handle auth buttons
  if (loginButton) {
    loginButton.addEventListener('click', () => showAuthModal('login'));
  }
  
  if (registerButton) {
    registerButton.addEventListener('click', () => showAuthModal('register'));
  }
  
  // Handle form switching
  if (switchToLogin) {
    switchToLogin.addEventListener('click', (e) => {
      e.preventDefault();
      switchAuthForm('login');
    });
  }
  
  if (switchToRegister) {
    switchToRegister.addEventListener('click', (e) => {
      e.preventDefault();
      switchAuthForm('register');
    });
  }
  
  window.closeAuthModal = () => {
    modal.classList.add('hidden');
    clearAuthError();
  };
  
  modal.addEventListener('click', (e) => {
    if (e.target === modal) {
      closeAuthModal();
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

function showAuthModal(mode = 'register') {
  const modal = document.getElementById('authModal');
  modal.classList.remove('hidden');
  switchAuthForm(mode);
}

function switchAuthForm(mode) {
  const registrationForm = document.getElementById('registrationForm');
  const loginForm = document.getElementById('loginForm');
  const modalTitle = document.getElementById('modalTitle');
  
  if (mode === 'login') {
    registrationForm.classList.add('hidden');
    loginForm.classList.remove('hidden');
    modalTitle.textContent = '🔑 Login to Your Account';
  } else {
    loginForm.classList.add('hidden');
    registrationForm.classList.remove('hidden');
    modalTitle.textContent = '🔒 Register to Unlock Premium Characters';
  }
  
  clearAuthError();
}

export function getSelectedCharacter() {
  return selectedCharacter;
}

// Authentication functions
async function checkAuthStatus() {
  try {
    const response = await fetch('/api/auth/me');
    const data = await response.json();
    
    if (data.success && data.authenticated) {
      currentUser = data.user;
      updateUIForLoggedInUser(data.user);
    } else {
      currentUser = null;
      updateUIForLoggedOutUser();
    }
  } catch (error) {
    console.error('Error checking auth status:', error);
    currentUser = null;
    updateUIForLoggedOutUser();
  }
}

function updateUIForLoggedInUser(user) {
  const userPanel = document.getElementById('userPanel');
  const loggedUsername = document.getElementById('loggedUsername');
  const authButtons = document.getElementById('authButtons');
  
  if (userPanel && loggedUsername) {
    loggedUsername.textContent = user.username;
    userPanel.classList.remove('hidden');
  }
  
  if (authButtons) {
    authButtons.classList.add('hidden');
  }
  
  // Unlock premium characters for registered users
  if (user.is_registered) {
    unlockAllCharacters();
  }
}

function updateUIForLoggedOutUser() {
  const userPanel = document.getElementById('userPanel');
  const authButtons = document.getElementById('authButtons');
  
  if (userPanel) {
    userPanel.classList.add('hidden');
  }
  
  if (authButtons) {
    authButtons.classList.remove('hidden');
  }
  
  // Keep premium characters locked
  lockPremiumCharacters();
}

function unlockAllCharacters() {
  const lockedBoxes = document.querySelectorAll('.character-box.locked');
  lockedBoxes.forEach(box => {
    box.classList.remove('locked');
    box.classList.add('unlocked');
    const lockIcon = box.querySelector('.lock-icon');
    if (lockIcon) {
      lockIcon.remove();
    }
  });
}

function lockPremiumCharacters() {
  const goldenBox = document.querySelector('[data-character="golden"]');
  const blackBox = document.querySelector('[data-character="black"]');
  
  [goldenBox, blackBox].forEach(box => {
    if (box) {
      box.classList.remove('unlocked', 'selected');
      box.classList.add('locked');
      if (!box.querySelector('.lock-icon')) {
        const lockIcon = document.createElement('span');
        lockIcon.className = 'lock-icon';
        lockIcon.textContent = '🔒';
        box.appendChild(lockIcon);
      }
    }
  });
  
  // Reset selection to default if premium character was selected
  if (selectedCharacter === 'golden' || selectedCharacter === 'black') {
    selectCharacter('boba');
  }
}

async function handleRegistration(e) {
  e.preventDefault();
  
  const formData = new FormData(e.target);
  const registrationData = {
    username: formData.get('username'),
    email: formData.get('email'),
    password: formData.get('password')
  };
  
  try {
    const response = await fetch('/api/auth/register', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(registrationData)
    });
    
    const data = await response.json();
    
    if (data.success) {
      currentUser = data.user;
      updateUIForLoggedInUser(data.user);
      closeAuthModal();
      
      // Show success message
      alert('🎉 Registration successful! Premium characters unlocked!');
    } else {
      showAuthError(data.error || 'Registration failed');
    }
  } catch (error) {
    console.error('Registration error:', error);
    showAuthError('Network error. Please try again.');
  }
}

async function handleLogin(e) {
  e.preventDefault();
  
  const formData = new FormData(e.target);
  const loginData = {
    username: formData.get('username'),
    password: formData.get('password')
  };
  
  try {
    const response = await fetch('/api/auth/login', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(loginData)
    });
    
    const data = await response.json();
    
    if (data.success) {
      currentUser = data.user;
      updateUIForLoggedInUser(data.user);
      closeAuthModal();
      
      // Show success message
      alert('🎉 Login successful! Welcome back!');
    } else {
      showAuthError(data.error || 'Login failed');
    }
  } catch (error) {
    console.error('Login error:', error);
    showAuthError('Network error. Please try again.');
  }
}

async function handleLogout() {
  try {
    const response = await fetch('/api/auth/logout', {
      method: 'POST',
    });
    
    const data = await response.json();
    
    if (data.success) {
      currentUser = null;
      updateUIForLoggedOutUser();
      alert('👋 Logged out successfully!');
    }
  } catch (error) {
    console.error('Logout error:', error);
  }
}

function showAuthError(message) {
  const errorDiv = document.getElementById('authError');
  if (errorDiv) {
    errorDiv.textContent = message;
    errorDiv.classList.remove('hidden');
  }
}

function clearAuthError() {
  const errorDiv = document.getElementById('authError');
  if (errorDiv) {
    errorDiv.classList.add('hidden');
    errorDiv.textContent = '';
  }
}