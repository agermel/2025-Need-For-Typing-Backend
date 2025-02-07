//go:build wireinject
// +build wireinject

package routes

import (
	"type/controllers"
	"type/dao"
	"type/service"

	"github.com/google/wire"
)

func InitSongController() *controllers.SongController {
	wire.Build(
		dao.NewSongDAO,
		service.NewSongService,
		controllers.NewSongController,
	)
	return &controllers.SongController{}
}

func InitScoreController() *controllers.ScoreController {
	wire.Build(
		dao.NewScoreDAO,
		service.NewScoreService,
		controllers.NewScoreController,
	)
	return &controllers.ScoreController{}
}

func InitUserController() *controllers.UserController {
	wire.Build(
		dao.NewUserDAO,
		service.NewUserService,
		controllers.NewUserController,
	)
	return &controllers.UserController{}
}

func InitAssetController() *controllers.AssetController {
	wire.Build(
		dao.NewAssetDAO,
		service.NewAssetService,
		controllers.NewAssetController,
	)
	return &controllers.AssetController{}
}
