package chat

import (
	"encoding/json"
	"sync"
)

// ChatRoom 聊天室結構
type ChatRoom struct {
	clients    map[*Client]bool
	Broadcast  chan []byte
	Register   chan *Client
	Unregister chan *Client
	messages   []Message
	mutex      sync.RWMutex
}

// NewChatRoom 建立新的聊天室
func NewChatRoom() *ChatRoom {
	return &ChatRoom{
		clients:    make(map[*Client]bool),
		Broadcast:  make(chan []byte),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		messages:   make([]Message, 0),
	}
}

// Run 運行聊天室
func (cr *ChatRoom) Run() {
	for {
		select {
		case client := <-cr.Register:
			cr.clients[client] = true
			// 發送加入通知
			joinMsg := NewJoinMessage(client.Username)
			cr.broadcastMessage(joinMsg)
			// 發送歷史訊息
			cr.sendHistoryToClient(client)

		case client := <-cr.Unregister:
			if _, ok := cr.clients[client]; ok {
				delete(cr.clients, client)
				client.Close()
				// 發送離開通知
				leaveMsg := NewLeaveMessage(client.Username)
				cr.broadcastMessage(leaveMsg)
			}

		case message := <-cr.Broadcast:
			// 直接廣播訊息給所有客戶端
			for client := range cr.clients {
				select {
				case client.Send <- message:
				default:
					client.Close()
					delete(cr.clients, client)
				}
			}
		}
	}
}

// broadcastMessage 廣播訊息到所有客戶端
func (cr *ChatRoom) broadcastMessage(msg Message) {
	data, _ := json.Marshal(msg)
	for client := range cr.clients {
		select {
		case client.Send <- data:
		default:
			client.Close()
			delete(cr.clients, client)
		}
	}
}

// sendHistoryToClient 發送歷史訊息給新客戶端
func (cr *ChatRoom) sendHistoryToClient(client *Client) {
	cr.mutex.RLock()
	defer cr.mutex.RUnlock()

	for _, msg := range cr.messages {
		data, _ := json.Marshal(msg)
		select {
		case client.Send <- data:
		default:
			return
		}
	}
}

// AddMessage 添加新訊息到聊天室
func (cr *ChatRoom) AddMessage(msg Message) {
	cr.mutex.Lock()
	defer cr.mutex.Unlock()
	
	cr.messages = append(cr.messages, msg)
	if len(cr.messages) > 100 {
		cr.messages = cr.messages[1:]
	}
}

// GetMessages 獲取所有訊息
func (cr *ChatRoom) GetMessages() []Message {
	cr.mutex.RLock()
	defer cr.mutex.RUnlock()
	
	messages := make([]Message, len(cr.messages))
	copy(messages, cr.messages)
	return messages
}

// GetClientCount 獲取客戶端數量
func (cr *ChatRoom) GetClientCount() int {
	cr.mutex.RLock()
	defer cr.mutex.RUnlock()
	return len(cr.clients)
}

// GetMessageCount 獲取訊息數量
func (cr *ChatRoom) GetMessageCount() int {
	cr.mutex.RLock()
	defer cr.mutex.RUnlock()
	return len(cr.messages)
}
