# 多階段建構 Dockerfile
FROM golang:1.24.5-alpine AS builder

# 設定工作目錄
WORKDIR /app

# 安裝必要的套件
RUN apk add --no-cache git

# 複製 go mod 檔案
COPY go.mod go.sum ./

# 下載依賴
RUN go mod download

# 複製原始碼
COPY . .

# 建構應用程式
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# 執行階段
FROM alpine:latest

# 安裝 ca-certificates 用於 HTTPS
RUN apk --no-cache add ca-certificates

# 建立非 root 用戶
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

# 設定工作目錄
WORKDIR /app

# 從 builder 階段複製二進制檔案
COPY --from=builder /app/main .

# 複製模板和靜態檔案
COPY --from=builder /app/templates ./templates
COPY --from=builder /app/static ./static

# 變更檔案擁有者
RUN chown -R appuser:appgroup /app

# 切換到非 root 用戶
USER appuser

# 暴露端口
EXPOSE 8080

# 健康檢查
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/api/stats || exit 1

# 啟動應用程式
CMD ["./main"]
