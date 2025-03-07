package service

import (
	"context"
	"time"
	"type/dao"
	"type/database"
	"type/models"
)

func SaveCode(ctx context.Context, email, code string, ttl time.Duration) error {
	err := dao.SaveCodeWithEmail(ctx, email, code, ttl, database.Rdb)
	if err != nil {
		return err
	}

	return nil
}

func (service *UserService) VerifyCode(ctx context.Context, email, code string) (*models.User, error) {
	err := dao.GetCodeByEmail(ctx, email, code, database.Rdb)
	if err != nil {
		return nil, err
	}

	user, err := service.UserDAO.GetUserByEmail(email)
	if err != nil {
		return nil, err
	}

	// 标记邮箱已验证
	err = service.UserDAO.VerifyEmail(email)
	if err != nil {
		return nil, err
	}

	err = dao.DeleteCode(ctx, email, database.Rdb)
	if err != nil {
		return nil, err
	}

	return user, nil
}
