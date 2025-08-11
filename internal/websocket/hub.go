package websocket

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go-chat-system/internal/chat"
	"go-chat-system/pkg/utils"
)

// Hub WebSocket 中心
type Hub struct {
	chatRoom *chat.ChatRoom
}

// NewHub 創建新的 WebSocket 中心
func NewHub(chatRoom *chat.ChatRoom) *Hub {
	return &Hub{
		chatRoom: chatRoom,
	}
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // 允許所有來源，生產環境應該限制
	},
}

// HandleWebSocket WebSocket 處理函數
func (h *Hub) HandleWebSocket(c *gin.Context) {
	username := c.Query("username")
	if username == "" {
		username = "Anonymous"
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("WebSocket upgrade failed: %v", err)
		return
	}

	client := chat.NewClient(utils.GenerateID(), username, conn)

	h.chatRoom.Register <- client

	go h.writePump(client)
	go h.readPump(client)
}

// readPump 讀取 WebSocket 訊息
func (h *Hub) readPump(client *chat.Client) {
	defer func() {
		h.chatRoom.Unregister <- client
		client.Close()
	}()

	for {
		_, message, err := client.Conn.ReadMessage()
		if err != nil {
			break
		}

		// 解析訊息
		var msg chat.Message
		if err := json.Unmarshal(message, &msg); err != nil {
			msg = chat.NewMessage(utils.GenerateID(), client.Username, string(message), "message")
		}

		// 儲存訊息
		h.chatRoom.AddMessage(msg)

		// 廣播結構化的訊息
		messageData, _ := json.Marshal(msg)
		h.chatRoom.Broadcast <- messageData
	}
}

// writePump 寫入 WebSocket 訊息
func (h *Hub) writePump(client *chat.Client) {
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
