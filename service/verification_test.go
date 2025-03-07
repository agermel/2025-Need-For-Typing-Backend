package service_test

import (
	"context"
	"testing"
	"time"
	"type/mocks"
	"type/models"
	"type/service"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestVerification(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// 创建 Mock 对象
	mockUserDAO := mocks.NewMockUserDAOInterface(ctrl)

	Mockservice := &service.UserService{
		UserDAO: mockUserDAO,
	}

	ctx := context.Background()
	email := "Aoi@Sone.com"
	code := "123456"
	user := &models.User{ID: 1, Username: "Aoi", Email: email}

	mockUserDAO.EXPECT().GetUserByEmail(email).Return(user, nil)
	mockUserDAO.EXPECT().VerifyEmail(email).Return(nil)

	// 虽说是userService的方法，但是只要给他指定是哪个(使用gomock生成的符合userservice接口的)service就不久好了吗
	service.SaveCode(ctx, email, code, 5*time.Second)
	result, err := Mockservice.VerifyCode(ctx, email, code)

	require.NoError(t, err)
	require.NotNil(t, result)
	require.Equal(t, user.Username, result.Username)

}

func BenchmarkVerification(b *testing.B) {
	ctrl := gomock.NewController(b)
	defer ctrl.Finish()

	// 创建 Mock 对象
	mockUserDAO := mocks.NewMockUserDAOInterface(ctrl)

	Mockservice := &service.UserService{
		UserDAO: mockUserDAO,
	}

	ctx := context.Background()
	email := "Aoi@Sone.com"
	code := "123456"
	user := &models.User{ID: 1, Username: "Aoi", Email: email}

	// .AnyTimes() 允许 Mock 对象多次被调用，避免 gomock 限制只调用一次。
	mockUserDAO.EXPECT().GetUserByEmail(email).Return(user, nil).AnyTimes()
	mockUserDAO.EXPECT().VerifyEmail(email).Return(nil).AnyTimes()

	b.ResetTimer() // 重置计时器，确保只测量 `b.N` 次调用的时间

	for i := 0; i < b.N; i++ {
		service.SaveCode(ctx, email, code, 5*time.Second)
		result, err := Mockservice.VerifyCode(ctx, email, code)

		require.NoError(b, err)
		require.NotNil(b, result)
		require.Equal(b, user.Username, result.Username)
	}
}
