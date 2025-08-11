// 聊天系統前端邏輯
let ws = null;
let currentUsername = '';
let isConnected = false;

// DOM 元素
const loginSection = document.getElementById('login-section');
const chatSection = document.getElementById('chat-section');
const usernameInput = document.getElementById('username-input');
const messageInput = document.getElementById('message-input');
const chatMessages = document.getElementById('chat-messages');
const onlineCount = document.getElementById('online-count');
const messageCount = document.getElementById('message-count');

// 加入聊天室
function joinChat() {
    const username = usernameInput.value.trim();
    if (!username) {
        alert('請輸入暱稱！');
        return;
    }
    
    currentUsername = username;
    connectWebSocket();
}

// 連接 WebSocket
function connectWebSocket() {
    const wsUrl = `ws://${window.location.host}/ws?username=${encodeURIComponent(currentUsername)}`;
    
    try {
        ws = new WebSocket(wsUrl);
        
        ws.onopen = function() {
            console.log('WebSocket 連接成功');
            isConnected = true;
            showChatSection();
            updateStats();
        };
        
        ws.onmessage = function(event) {
            handleMessage(event.data);
        };
        
        ws.onclose = function() {
            console.log('WebSocket 連接關閉');
            isConnected = false;
            showLoginSection();
        };
        
        ws.onerror = function(error) {
            console.error('WebSocket 錯誤:', error);
            alert('連接失敗，請重試！');
        };
        
    } catch (error) {
        console.error('建立 WebSocket 連接失敗:', error);
        alert('連接失敗，請重試！');
    }
}

// 處理接收到的訊息
function handleMessage(data) {
    try {
        const message = JSON.parse(data);
        displayMessage(message);
        updateStats();
    } catch (error) {
        console.error('解析訊息失敗:', error);
    }
}

// 顯示訊息
function displayMessage(message) {
    const messageDiv = document.createElement('div');
    messageDiv.className = `message ${getMessageClass(message)}`;
    
    const timestamp = new Date(message.timestamp).toLocaleTimeString('zh-TW');
    
    messageDiv.innerHTML = `
        <div class="message-header">${message.username}</div>
        <div class="message-content">${escapeHtml(message.content)}</div>
        <div class="message-time">${timestamp}</div>
    `;
    
    chatMessages.appendChild(messageDiv);
    chatMessages.scrollTop = chatMessages.scrollHeight;
}

// 獲取訊息樣式類別
function getMessageClass(message) {
    if (message.type === 'join' || message.type === 'leave') {
        return 'system';
    } else if (message.username === currentUsername) {
        return 'user';
    } else {
        return 'other';
    }
}

// 發送訊息
function sendMessage() {
    const content = messageInput.value.trim();
    if (!content || !isConnected) {
        return;
    }
    
    const message = {
        username: currentUsername,
        content: content,
        timestamp: new Date().toISOString(),
        type: 'message'
    };
    
    if (ws && ws.readyState === WebSocket.OPEN) {
        ws.send(JSON.stringify(message));
        messageInput.value = '';
    }
}

// 按 Enter 鍵發送訊息
messageInput.addEventListener('keypress', function(e) {
    if (e.key === 'Enter') {
        sendMessage();
    }
});

// 顯示聊天區域
function showChatSection() {
    loginSection.style.display = 'none';
    chatSection.style.display = 'block';
    messageInput.focus();
}

// 顯示登入區域
function showLoginSection() {
    chatSection.style.display = 'none';
    loginSection.style.display = 'block';
    currentUsername = '';
    ws = null;
}

// 更新統計資訊
async function updateStats() {
    try {
        const response = await fetch('/api/stats');
        const stats = await response.json();
        
        onlineCount.textContent = stats.online_users;
        messageCount.textContent = stats.total_messages;
    } catch (error) {
        console.error('獲取統計資訊失敗:', error);
    }
}

// 定期更新統計資訊
setInterval(updateStats, 5000);

// 頁面載入時更新統計
document.addEventListener('DOMContentLoaded', function() {
    updateStats();
});

// HTML 轉義函數
function escapeHtml(text) {
    const div = document.createElement('div');
    div.textContent = text;
    return div.innerHTML;
}

// 頁面卸載時關閉連接
window.addEventListener('beforeunload', function() {
    if (ws) {
        ws.close();
    }
});

// 連接狀態指示器
function updateConnectionStatus() {
    const statusIndicator = document.createElement('div');
    statusIndicator.id = 'connection-status';
    statusIndicator.style.cssText = `
        position: fixed;
        top: 20px;
        right: 20px;
        padding: 10px 15px;
        border-radius: 20px;
        color: white;
        font-weight: bold;
        z-index: 1000;
        transition: all 0.3s ease;
    `;
    
    if (isConnected) {
        statusIndicator.style.background = '#28a745';
        statusIndicator.textContent = '🟢 已連接';
    } else {
        statusIndicator.style.background = '#dc3545';
        statusIndicator.textContent = '🔴 未連接';
    }
    
    // 移除舊的狀態指示器
    const oldIndicator = document.getElementById('connection-status');
    if (oldIndicator) {
        oldIndicator.remove();
    }
    
    document.body.appendChild(statusIndicator);
    
    // 3秒後自動隱藏
    setTimeout(() => {
        if (statusIndicator.parentNode) {
            statusIndicator.style.opacity = '0';
            setTimeout(() => {
                if (statusIndicator.parentNode) {
                    statusIndicator.remove();
                }
            }, 300);
        }
    }, 3000);
}

// 監聽連接狀態變化
setInterval(updateConnectionStatus, 1000);
