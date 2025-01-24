package utils

import (
	"math/rand"
	"time"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// GenerateRoomID 生成随机五位房间号
func GenerateRoomID() string {
	rand.Seed(time.Now().UnixNano()) // 初始化随机种子
	roomID := make([]byte, 5)        // 创建长度为 5 的字节切片
	for i := range roomID {
		roomID[i] = charset[rand.Intn(len(charset))] // 从字符集中随机取值
	}
	return string(roomID) // 转换为字符串
}

// 生成随机8位验证代码
func GenerateRandomVerifyCode() string {
	rand.Seed(time.Now().UnixNano()) // 初始化随机种子
	VerifyCode := make([]byte, 8)    // 创建长度为 5 的字节切片
	for i := range VerifyCode {
		VerifyCode[i] = charset[rand.Intn(len(charset))] // 从字符集中随机取值
	}
	return string(VerifyCode) // 转换为字符串
}
