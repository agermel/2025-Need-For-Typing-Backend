package main

import (
	"type/config"
	"type/database"

	"github.com/gin-gonic/gin"
)

// @title NeedForTyping
// @version 1.0
// @description 一个打字游戏

func main() {
	config.LoadConfig()

	database.InitDatabase()
	defer database.CloseDatabase()
	app := InitApp()
	app.Run()
}

type App struct {
	e *gin.Engine
}

func NewApp(e *gin.Engine) App {
	return App{
		e: e,
	}
}
func (a *App) Run() {
	a.e.Run(":8080")
}
