package utils

import "time"

// GenerateID 生成簡單的 ID
func GenerateID() string {
	return time.Now().Format("20060102150405")
}
