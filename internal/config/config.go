package config

import (
	"os"
	"strconv"
)

// Config 應用配置
type Config struct {
	Port         string
	MaxMessages  int
	AllowOrigins []string
}

// Load 載入配置
func Load() *Config {
	port := getEnv("PORT", "8080")
	maxMessages, _ := strconv.Atoi(getEnv("MAX_MESSAGES", "100"))

	return &Config{
		Port:         port,
		MaxMessages:  maxMessages,
		AllowOrigins: []string{"*"}, // 生產環境應該限制
	}
}

// getEnv 獲取環境變數，如果不存在則返回預設值
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
