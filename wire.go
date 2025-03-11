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
		dao.NewGameDAO,
		service.NewSongService,
		service.NewScoreService,
		service.NewUserService,
		service.NewAssetService,
		service.NewGameService,
		controllers.NewAssetController,
		controllers.NewUserController,
		controllers.NewSongController,
		controllers.NewScoreController,
		controllers.NewGameController,
		routes.RegisterRoutes,
		NewApp,
	)
	return App{}
}
