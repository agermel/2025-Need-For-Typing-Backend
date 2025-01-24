package main

import (
	"type/config"
	"type/database"
	"type/middlewares"
	"type/routes"

	"github.com/gin-gonic/gin"
)

// @title NeedForTyping
// @version 1.0
// @description 一个打字游戏

func main() {
	config.LoadConfig()

	database.InitDatabase()
	defer database.CloseDatabase()

	r := gin.Default()

	r.Use(middlewares.CORSMiddleware())

	routes.RegisterRoutes(r)

	r.Run(":8888")
}
