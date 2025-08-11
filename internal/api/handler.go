package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go-chat-system/internal/chat"
)

// Handler API 處理器
type Handler struct {
	chatRoom *chat.ChatRoom
	startTime time.Time
}

// NewHandler 創建新的 API 處理器
func NewHandler(chatRoom *chat.ChatRoom) *Handler {
	return &Handler{
		chatRoom: chatRoom,
		startTime: time.Now(),
	}
}

// GetMessages 獲取聊天記錄
func (h *Handler) GetMessages(c *gin.Context) {
	messages := h.chatRoom.GetMessages()
	
	c.JSON(http.StatusOK, gin.H{
		"messages": messages,
		"count":    len(messages),
	})
}

// GetStats 獲取系統統計
func (h *Handler) GetStats(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"online_users":   h.chatRoom.GetClientCount(),
		"total_messages": h.chatRoom.GetMessageCount(),
		"uptime":         time.Since(h.startTime).String(),
	})
}
