//go:build wireinject
// +build wireinject

package main

import (
	"type/controllers"
	"type/dao"
	"type/routes"
	"type/service"

	"github.com/google/wire"
)

func InitApp() App {
	wire.Build(
		dao.NewScoreDAO,
		dao.NewAssetDAO,
		dao.NewUserDAO,
		dao.NewSongDAO,
		dao.NewQDAO,
		service.NewSongService,
		service.NewScoreService,
		service.NewUserService,
		service.NewAssetService,
		service.NewQuestionService,
		controllers.NewAssetController,
		controllers.NewUserController,
		controllers.NewSongController,
		controllers.NewScoreController,
		controllers.NewQuestionController,
		routes.RegisterRoutes,
		NewApp,
	)
	return App{}
}
