package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// Message 聊天訊息結構
type Message struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	Content   string    `json:"content"`
	Timestamp time.Time `json:"timestamp"`
	Type      string    `json:"type"` // "message", "join", "leave"
}

// Client WebSocket 客戶端結構
type Client struct {
	ID       string
	Username string
	Conn     *websocket.Conn
	Send     chan []byte
}

// ChatRoom 聊天室結構
type ChatRoom struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
	messages   []Message
	mutex      sync.RWMutex
}

// Hub 聊天系統中心
type Hub struct {
	chatRoom *ChatRoom
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // 允許所有來源，生產環境應該限制
	},
}

// NewChatRoom 建立新的聊天室
func NewChatRoom() *ChatRoom {
	return &ChatRoom{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		messages:   make([]Message, 0),
	}
}

// Run 運行聊天室
func (cr *ChatRoom) Run() {
	for {
		select {
		case client := <-cr.register:
			cr.clients[client] = true
			// 發送加入通知
			joinMsg := Message{
				ID:        generateID(),
				Username:  "System",
				Content:   client.Username + " 加入了聊天室",
				Timestamp: time.Now(),
				Type:      "join",
			}
			cr.broadcastMessage(joinMsg)
			// 發送歷史訊息
			cr.sendHistoryToClient(client)

		case client := <-cr.unregister:
			if _, ok := cr.clients[client]; ok {
				delete(cr.clients, client)
				close(client.Send)
				// 發送離開通知
				leaveMsg := Message{
					ID:        generateID(),
					Username:  "System",
					Content:   client.Username + " 離開了聊天室",
					Timestamp: time.Now(),
					Type:      "leave",
				}
				cr.broadcastMessage(leaveMsg)
			}

		case message := <-cr.broadcast:
			// 直接廣播訊息給所有客戶端，不需要重複儲存
			for client := range cr.clients {
				select {
				case client.Send <- message:
				default:
					close(client.Send)
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
			close(client.Send)
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

// generateID 生成簡單的 ID
func generateID() string {
	return time.Now().Format("20060102150405")
}

// handleWebSocket WebSocket 處理函數
func (h *Hub) handleWebSocket(c *gin.Context) {
	username := c.Query("username")
	if username == "" {
		username = "Anonymous"
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("WebSocket upgrade failed: %v", err)
		return
	}

	client := &Client{
		ID:       generateID(),
		Username: username,
		Conn:     conn,
		Send:     make(chan []byte, 256),
	}

	h.chatRoom.register <- client

	go h.writePump(client)
	go h.readPump(client)
}

// readPump 讀取 WebSocket 訊息
func (h *Hub) readPump(client *Client) {
	defer func() {
		h.chatRoom.unregister <- client
		client.Conn.Close()
	}()

	for {
		_, message, err := client.Conn.ReadMessage()
		if err != nil {
			break
		}

		// 解析訊息
		var msg Message
		if err := json.Unmarshal(message, &msg); err != nil {
			msg = Message{
				ID:        generateID(),
				Username:  client.Username,
				Content:   string(message),
				Timestamp: time.Now(),
				Type:      "message",
			}
		}

		// 儲存訊息
		h.chatRoom.mutex.Lock()
		h.chatRoom.messages = append(h.chatRoom.messages, msg)
		if len(h.chatRoom.messages) > 100 {
			h.chatRoom.messages = h.chatRoom.messages[1:]
		}
		h.chatRoom.mutex.Unlock()

		// 廣播結構化的訊息，而不是原始訊息
		messageData, _ := json.Marshal(msg)
		h.chatRoom.broadcast <- messageData
	}
}

// writePump 寫入 WebSocket 訊息
func (h *Hub) writePump(client *Client) {
	defer client.Conn.Close()

	for message := range client.Send {
		w, err := client.Conn.NextWriter(websocket.TextMessage)
		if err != nil {
			return
		}
		w.Write(message)

		if err := w.Close(); err != nil {
			return
		}
	}
}

// getMessages 獲取聊天記錄
func (h *Hub) getMessages(c *gin.Context) {
	h.chatRoom.mutex.RLock()
	defer h.chatRoom.mutex.RUnlock()

	c.JSON(http.StatusOK, gin.H{
		"messages": h.chatRoom.messages,
		"count":    len(h.chatRoom.messages),
	})
}

// getStats 獲取系統統計
func (h *Hub) getStats(c *gin.Context) {
	h.chatRoom.mutex.RLock()
	defer h.chatRoom.mutex.RUnlock()

	c.JSON(http.StatusOK, gin.H{
		"online_users":   len(h.chatRoom.clients),
		"total_messages": len(h.chatRoom.messages),
		"uptime":         time.Since(startTime).String(),
	})
}

var startTime = time.Now()

func main() {
	// 建立聊天室
	chatRoom := NewChatRoom()
	hub := &Hub{chatRoom: chatRoom}

	// 啟動聊天室
	go chatRoom.Run()

	// 設定 Gin 路由
	r := gin.Default()

	// CORS 設定
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// API 路由
	api := r.Group("/api")
	{
		api.GET("/messages", hub.getMessages)
		api.GET("/stats", hub.getStats)
	}

	// WebSocket 路由
	r.GET("/ws", hub.handleWebSocket)

	// 靜態檔案服務
	r.Static("/static", "./static")
	r.LoadHTMLGlob("templates/*")

	// 首頁
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	log.Println("聊天系統啟動在 :8080 端口")
	log.Fatal(r.Run(":8080"))
}
