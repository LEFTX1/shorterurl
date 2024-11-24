package testutil

import (
	"crypto/rand"
	"encoding/hex"
	"time"
)

// GenerateTestUsername 生成一个随机的测试用户名
func GenerateTestUsername() string {
	bytes := make([]byte, 16)                   // 创建一个长度为16的字节切片
	if _, err := rand.Read(bytes); err != nil { // 生成随机字节
		return "test_" + hex.EncodeToString([]byte(time.Now().String())) // 如果生成随机字节失败，使用当前时间生成用户名
	}
	return "test_" + hex.EncodeToString(bytes) // 返回生成的随机用户名
}
