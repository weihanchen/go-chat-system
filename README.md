# 💬 Go 聊天系統

一個使用 Go 語言建構的即時線上聊天系統，展示後端開發技能。

## 🚀 專案特色

- **即時通訊**: 使用 WebSocket 實現即時聊天功能
- **無資料庫依賴**: 使用記憶體儲存聊天記錄和用戶狀態
- **並發處理**: 展示 Go 的 goroutine 和 channel 使用
- **現代化 UI**: 響應式設計，支援桌面和行動裝置
- **RESTful API**: 提供聊天記錄和系統統計 API
- **CORS 支援**: 跨域請求支援

## 🛠️ 技術架構

### 後端技術
- **Go 1.21+**: 主要程式語言
- **Gin**: 輕量級 Web 框架
- **Gorilla WebSocket**: WebSocket 實現
- **Gin CORS**: 跨域資源共享支援

### 前端技術
- **原生 JavaScript**: 無框架依賴
- **WebSocket API**: 即時通訊
- **CSS3**: 現代化樣式和動畫
- **響應式設計**: 支援多種裝置

## 📁 專案結構

```
go-chat-system/
├── main.go              # 主要程式檔案
├── go.mod               # Go 模組依賴
├── go.sum               # 依賴校驗檔案
├── README.md            # 專案說明
├── templates/           # HTML 模板
│   └── index.html      # 主頁面
└── static/             # 靜態資源
    ├── style.css       # 樣式檔案
    └── script.js       # JavaScript 邏輯
```

## 🚀 快速開始

### 前置需求
- Go 1.24.5 或更新版本
- 現代化瀏覽器（支援 WebSocket）
- Docker（可選，用於容器化部署）

### 方法一：直接執行（開發環境）

1. **克隆專案**
   ```bash
   git clone <repository-url>
   cd go-chat-system
   ```

2. **安裝依賴**
   ```bash
   go mod tidy
   ```

3. **執行專案**
   ```bash
   go run main.go
   ```

4. **開啟瀏覽器**
   訪問 `http://localhost:8080`

### 方法二：Docker 部署（生產環境）

1. **建構 Docker 映像檔**
   ```bash
   docker build -t go-chat-system:latest .
   ```

2. **使用 Docker 執行**
   ```bash
   docker run -p 8080:8080 go-chat-system:latest
   ```

3. **使用 Docker Compose（推薦）**
   ```bash
   docker compose up -d
   ```

4. **開啟瀏覽器**
   訪問 `http://localhost:8080`

## 🔧 功能說明

### 核心功能
- **用戶登入**: 輸入暱稱加入聊天室
- **即時聊天**: 支援多人即時對話
- **系統通知**: 用戶加入/離開通知
- **聊天記錄**: 顯示最近 100 條訊息
- **線上統計**: 即時顯示線上人數和訊息數量

### API 端點
- `GET /api/messages`: 獲取聊天記錄
- `GET /api/stats`: 獲取系統統計
- `GET /ws`: WebSocket 連接端點

### WebSocket 事件
- **連接**: 自動發送歷史訊息
- **訊息**: 即時廣播給所有用戶
- **斷線**: 自動清理用戶狀態

## 💡 後端技能展示

### 1. 並發處理 (Multi-threading)
```go
// 使用 goroutine 處理多個客戶端
go h.writePump(client)
go h.readPump(client)

// 使用 channel 進行協程間通訊
case client := <-cr.register:
case message := <-cr.broadcast:
```

### 2. 容器化部署 (Docker)
```dockerfile
# 多階段建構，優化映像檔大小
FROM golang:1.21-alpine AS builder
# ... 建構階段 ...

FROM alpine:latest
# ... 執行階段 ...
```

### 2. 記憶體管理
```go
// 限制聊天記錄數量，防止記憶體洩漏
if len(cr.messages) > 100 {
    cr.messages = cr.messages[1:]
}
```

### 3. 同步機制
```go
// 使用 RWMutex 保護共享資源
cr.mutex.RLock()
defer cr.mutex.RUnlock()
```

### 4. 錯誤處理
```go
// 優雅的錯誤處理和資源清理
defer func() {
    h.chatRoom.unregister <- client
    client.Conn.Close()
}()
```

### 5. 網路程式設計
```go
// WebSocket 升級和處理
conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
if err != nil {
    log.Printf("WebSocket upgrade failed: %v", err)
    return
}
```

## 🌟 專案亮點

### 架構設計
- **模組化設計**: 清晰的結構分離
- **可擴展性**: 易於添加新功能
- **錯誤處理**: 完善的錯誤處理機制

### 效能優化
- **記憶體管理**: 自動清理過期資料
- **連接池**: 高效的客戶端管理
- **非阻塞 I/O**: 使用 channel 避免阻塞

### 用戶體驗
- **即時反饋**: 連接狀態指示器
- **響應式設計**: 支援多種裝置
- **直觀介面**: 簡潔美觀的使用者介面

## 🔍 測試建議

### 功能測試
1. 開啟多個瀏覽器視窗
2. 測試用戶加入/離開通知
3. 驗證即時訊息傳遞
4. 檢查聊天記錄保存

### 壓力測試
1. 模擬大量用戶同時連接
2. 測試長時間運行的穩定性
3. 驗證記憶體使用情況

## 🚧 未來改進

- [ ] 添加用戶認證系統
- [ ] 支援私聊功能
- [ ] 添加檔案傳輸
- [ ] 實現聊天室分組
- [ ] 添加訊息搜尋功能
- [ ] 支援表情符號和圖片
