package main

import (
	"log"

	"go-chat-system/internal/api"
	"go-chat-system/internal/chat"
	"go-chat-system/internal/config"
	"go-chat-system/internal/websocket"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// 載入配置
	cfg := config.Load()

	// 建立聊天室
	chatRoom := chat.NewChatRoom()

	// 啟動聊天室
	go chatRoom.Run()

	// 建立 WebSocket 中心
	hub := websocket.NewHub(chatRoom)

	// 建立 API 處理器
	apiHandler := api.NewHandler(chatRoom)

	// 設定 Gin 路由
	r := gin.Default()

	// CORS 設定
	r.Use(cors.New(cors.Config{
		AllowOrigins:     cfg.AllowOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// API 路由
	apiGroup := r.Group("/api")
	{
		apiGroup.GET("/messages", apiHandler.GetMessages)
		apiGroup.GET("/stats", apiHandler.GetStats)
	}

	// WebSocket 路由
	r.GET("/ws", hub.HandleWebSocket)

	// 靜態檔案服務
	r.Static("/static", "./web/static")
	r.LoadHTMLGlob("web/templates/*")

	// 首頁
	r.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", nil)
	})

	log.Printf("聊天系統啟動在 :%s 端口", cfg.Port)
	log.Fatal(r.Run(":" + cfg.Port))
}
