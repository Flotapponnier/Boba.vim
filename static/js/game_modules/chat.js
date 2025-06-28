let chatHistory = [];
let chatHistoryVisible = false;

export function initializeChatHistory() {
  document.addEventListener("keydown", function (event) {
    if (event.key === "\\") {
      toggleChatHistory();
      event.preventDefault();
    }
  });
}

export function addToChatHistory(message) {
  const timestamp = new Date().toLocaleTimeString();
  chatHistory.push({
    time: timestamp,
    message: message,
  });

  if (chatHistory.length > 50) {
    chatHistory.shift();
  }
}

export function toggleChatHistory() {
  chatHistoryVisible = !chatHistoryVisible;
  let chatWindow = document.getElementById("chatHistoryWindow");

  if (chatHistoryVisible) {
    if (!chatWindow) {
      createChatHistoryWindow();
    } else {
      chatWindow.style.display = "block";
      updateChatHistoryContent();
    }
  } else {
    if (chatWindow) {
      chatWindow.style.display = "none";
    }
  }
}

function createChatHistoryWindow() {
  const chatWindow = document.createElement("div");
  chatWindow.id = "chatHistoryWindow";
  chatWindow.className = "chat-history-window";

  chatWindow.innerHTML = `
    <div class="chat-history-header">
      <h3>Chat History</h3>
      <button class="chat-close-btn" onclick="window.chatModule.toggleChatHistory()">×</button>
    </div>
    <div class="chat-history-content" id="chatHistoryContent">
      ${chatHistory.length === 0 ? '<p class="no-history">No commands yet. Start moving with HJKL!</p>' : ""}
    </div>
    <div class="chat-history-footer">
      Press \\ to close | Total commands: <span id="commandCount">${chatHistory.length}</span>
    </div>
  `;

  document.body.appendChild(chatWindow);
  updateChatHistoryContent();
}

function updateChatHistoryContent() {
  const content = document.getElementById("chatHistoryContent");
  const commandCount = document.getElementById("commandCount");

  if (chatHistory.length === 0) {
    content.innerHTML =
      '<p class="no-history">No commands yet. Start moving with HJKL!</p>';
  } else {
    content.innerHTML = chatHistory
      .map(
        (entry) =>
          `<div class="chat-entry">
        <span class="chat-time">[${entry.time}]</span>
        <span class="chat-message">${entry.message}</span>
      </div>`,
      )
      .join("");

    content.scrollTop = content.scrollHeight;
  }

  if (commandCount) {
    commandCount.textContent = chatHistory.length;
  }
}
