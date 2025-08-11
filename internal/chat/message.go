package chat

import "time"

// Message 聊天訊息結構
type Message struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	Content   string    `json:"content"`
	Timestamp time.Time `json:"timestamp"`
	Type      string    `json:"type"` // "message", "join", "leave"
}

// NewMessage 創建新訊息
func NewMessage(id, username, content, msgType string) Message {
	return Message{
		ID:        id,
		Username:  username,
		Content:   content,
		Timestamp: time.Now(),
		Type:      msgType,
	}
}

// NewSystemMessage 創建系統訊息
func NewSystemMessage(content string) Message {
	return NewMessage("", "System", content, "system")
}

// NewJoinMessage 創建加入訊息
func NewJoinMessage(username string) Message {
	return NewMessage("", "System", username+" 加入了聊天室", "join")
}

// NewLeaveMessage 創建離開訊息
func NewLeaveMessage(username string) Message {
	return NewMessage("", "System", username+" 離開了聊天室", "leave")
}
