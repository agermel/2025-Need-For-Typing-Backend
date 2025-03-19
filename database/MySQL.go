package database

import (
	"fmt"
	"log"

	"type/config"
	"type/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDatabase() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.AllConfig.Database.User,
		config.AllConfig.Database.Password,
		config.AllConfig.Server.Address,
		config.AllConfig.Server.Port,
		config.AllConfig.Database.Name,
	)

	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	DB.AutoMigrate(&models.User{}, &models.TotalScore{}, &models.Song{}, &models.Asset{}, &models.Question{})
}

func CloseDatabase() {
	db, err := DB.DB()
	if err != nil {
		log.Fatalf("获取数据库实例失败: %v", err)
	}
	db.Close()
}
