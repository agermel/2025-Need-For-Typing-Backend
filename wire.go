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
		service.NewSongService,
		service.NewScoreService,
		service.NewUserService,
		service.NewAssetService,
		controllers.NewAssetController,
		controllers.NewUserController,
		controllers.NewSongController,
		controllers.NewScoreController,
		routes.RegisterRoutes,
		NewApp,
	)
	return App{}
}
