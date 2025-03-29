// Chat Support System

// Generate a unique session ID for this chat
const sessionId = 'chat_' + Math.random().toString(36).substring(2, 15);
let userName = 'Guest';
let userEmail = '';

// DOM elements
const chatButton = document.getElementById('chat-button');
const chatContainer = document.getElementById('chat-container');
const chatClose = document.getElementById('chat-close');
const chatMessages = document.getElementById('chat-messages');
const chatInput = document.getElementById('chat-input');
const chatSend = document.getElementById('chat-send');

// Базовый URL API
const API_BASE_URL = "http://localhost:5000"; // Измените на ваш URL сервера
let lastCheckTime = Date.now(); // Время последней проверки сообщений

// Initialize chat
function initChat() {
    // First get user info if available
    if (localStorage.getItem('userName')) {
        userName = localStorage.getItem('userName');
    }
    if (localStorage.getItem('userEmail')) {
        userEmail = localStorage.getItem('userEmail');
    }
    
    // Event listeners
    chatButton.addEventListener('click', toggleChat);
    chatClose.addEventListener('click', toggleChat);
    chatSend.addEventListener('click', sendMessage);
    chatInput.addEventListener('keypress', function(e) {
        if (e.key === 'Enter') {
            sendMessage();
        }
    });
    
    // Check for existing messages in this session
    loadChatHistory();
    
    // Initial connection to backend
    connectChat();
    
    // Добавим отладочную информацию
    console.log("Chat initialized with sessionId:", sessionId);
}

function toggleChat() {
    if (chatContainer.style.display === 'flex') {
        chatContainer.style.display = 'none';
    } else {
        chatContainer.style.display = 'flex';
        chatInput.focus();
        markMessagesAsRead();
    }
}

function connectChat() {
    // Let backend know a new chat session has started
    fetch(`${API_BASE_URL}/api/chat/connect`, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({
            sessionId: sessionId,
            userName: userName,
            userEmail: userEmail
        })
    })
    .then(response => response.json())
    .then(data => {
        if (data.success) {
            console.log('Chat session connected');
            // Start polling for new messages
            startMessagePolling();
        } else {
            console.error('Failed to connect chat session:', data);
        }
    })
    .catch(error => {
        console.error('Error connecting chat:', error);
    });
}

function sendMessage() {
    const messageText = chatInput.value.trim();
    if (!messageText) return;
    
    // Add message to chat interface
    addMessageToChat(messageText, 'user');
    
    // Clear input field
    chatInput.value = '';
    
    // Send message to server
    fetch(`${API_BASE_URL}/api/chat/message`, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({
            sessionId: sessionId,
            from: 'user',
            text: messageText,
            userName: userName,
            userEmail: userEmail
        })
    })
    .then(response => response.json())
    .then(data => {
        if (!data.success) {
            console.error('Error sending message:', data.error);
        } else {
            console.log('Message sent successfully:', data);
        }
    })
    .catch(error => {
        console.error('Error sending message:', error);
    });
}

function addMessageToChat(text, sender) {
    const messageElement = document.createElement('div');
    messageElement.classList.add('chat-message');
    messageElement.classList.add(sender === 'user' ? 'user-message' : 'support-message');
    messageElement.textContent = text;
    
    chatMessages.appendChild(messageElement);
    
    // Scroll to bottom
    chatMessages.scrollTop = chatMessages.scrollHeight;
    
    // Save to local storage
    saveChatHistory();
}

function startMessagePolling() {
    // Poll for new messages every 3 seconds
    setInterval(() => {
        checkNewMessages();
    }, 3000);
}

function checkNewMessages() {
    fetch(`${API_BASE_URL}/api/chat/messages?sessionId=${sessionId}&lastCheck=${lastCheckTime}`)
        .then(response => {
            if (!response.ok) {
                throw new Error(`HTTP error! Status: ${response.status}`);
            }
            return response.json();
        })
        .then(data => {
            console.log('Polling response:', data); // Отладка
            
            if (data.success && data.messages && data.messages.length > 0) {
                data.messages.forEach(message => {
                    // Проверяем, что это сообщение от поддержки и оно новое
                    if (message.from === 'support') {
                        console.log('Support message received:', message);
                        addMessageToChat(message.text, 'support');
                        
                        // If chat is not open, show notification
                        if (chatContainer.style.display !== 'flex') {
                            showChatNotification();
                        }
                    }
                });
                
                // Обновляем время последней проверки
                lastCheckTime = Date.now();
                
                // Mark messages as read if chat is open
                if (chatContainer.style.display === 'flex') {
                    markMessagesAsRead();
                    
                    // Send read status to server
                    fetch(`${API_BASE_URL}/api/chat/mark-read`, {
                        method: 'POST',
                        headers: {
                            'Content-Type': 'application/json'
                        },
                        body: JSON.stringify({
                            sessionId: sessionId
                        })
                    }).catch(err => console.error('Error marking messages as read:', err));
                }
            }
        })
        .catch(error => {
            console.error('Error polling messages:', error);
        });
}

function showChatNotification() {
    // Add notification indicator to chat button
    chatButton.classList.add('has-notification');
    
    // You could also play a sound or show a browser notification here
    console.log('New message notification shown');
}

function markMessagesAsRead() {
    chatButton.classList.remove('has-notification');
}

function saveChatHistory() {
    localStorage.setItem(`chatHistory_${sessionId}`, chatMessages.innerHTML);
}

function loadChatHistory() {
    const history = localStorage.getItem(`chatHistory_${sessionId}`);
    if (history) {
        chatMessages.innerHTML = history;
        chatMessages.scrollTop = chatMessages.scrollHeight;
    }
}

// Initialize chat when DOM is loaded
document.addEventListener('DOMContentLoaded', initChat);
