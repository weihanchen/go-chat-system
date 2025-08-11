package chat

import (
	"sync"

	"github.com/gorilla/websocket"
)

// Client WebSocket 客戶端結構
type Client struct {
	ID       string
	Username string
	Conn     *websocket.Conn
	Send     chan []byte
	closed   bool
	mutex    sync.Mutex
}

// NewClient 創建新客戶端
func NewClient(id, username string, conn *websocket.Conn) *Client {
	return &Client{
		ID:       id,
		Username: username,
		Conn:     conn,
		Send:     make(chan []byte, 256),
	}
}

// Close 關閉客戶端連接
func (c *Client) Close() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if !c.closed {
		c.closed = true
		close(c.Send)
		c.Conn.Close()
	}
}
