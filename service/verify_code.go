package service

import (
	"errors"
	"sync"
	"time"
	"type/models"
)

type VerificaitionCode struct {
	Code      string
	ExpiresAt time.Time
}

// 线程安全的全局存储
// 匿名结构体可以直接结合其他类型（如 sync.Mutex）构建全局变量，减少额外的类型声明。
var verificationStore = struct {
	sync.Mutex
	data map[string]VerificaitionCode
}{data: make(map[string]VerificaitionCode)}

// 写入redis?(feature)
// ttl : time to live
func SaveCode(email, code string, ttl time.Duration) {
	verificationStore.Lock()
	defer verificationStore.Unlock()
	// 锁住了 verificationStore 这个结构体中的 sync.Mutex
	// 从而防止多个 Goroutine 同时对 verificationStore.data 进行操作

	verificationStore.data[email] = VerificaitionCode{
		Code:      code,
		ExpiresAt: time.Now().Add(ttl),
	}
}

func (service *UserService) VerifyCode(email, code string) (*models.User, error) {
	verificationStore.Lock()
	defer verificationStore.Unlock()

	v, exists := verificationStore.data[email]
	if !exists || time.Now().After(v.ExpiresAt) {
		return nil, errors.New("code过期或不存在") // code过期了或未找到
	}
	user, err := service.userDAO.GetUserByEmail(email)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// DeleteCode removes a verification code after successful verification
func DeleteCode(email string) {
	verificationStore.Lock()
	defer verificationStore.Unlock()

	delete(verificationStore.data, email)
}
