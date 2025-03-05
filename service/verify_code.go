package service

import (
	"context"
	"time"
	"type/dao"
	"type/database"
	"type/models"
)

var rdb = database.NewRedis()

func SaveCode(ctx context.Context, email, code string, ttl time.Duration) error {
	err := dao.SaveCodeWithEmail(ctx, email, code, ttl, database.Rdb)
	if err != nil {
		return err
	}

	return nil
}

func (service *UserService) VerifyCode(ctx context.Context, email, code string) (*models.User, error) {
	err := dao.GetCodeByEmail(ctx, email, code, rdb)
	if err != nil {
		return nil, err
	}

	user, err := service.userDAO.GetUserByEmail(email)
	if err != nil {
		return nil, err
	}

	// 标记邮箱已验证
	err = service.userDAO.VerifyEmail(email)
	if err != nil {
		return nil, err
	}

	err = dao.DeleteCode(ctx, email, rdb)
	if err != nil {
		return nil, err
	}

	return user, nil
}
