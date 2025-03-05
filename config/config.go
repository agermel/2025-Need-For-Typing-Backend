package config

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

// 配置结构体
type Config struct {
	Server struct {
		Port     int    `mapstructure:"port"`
		Address  string `mapstructure:"address"`
		RedisURL string `mapstructure:"redisurl"`
	} `mapstructure:"server"`

	Database struct {
		User          string `mapstructure:"user"`
		Password      string `mapstructure:"password"`
		Name          string `mapstructure:"name"`
		RedisPassword string `mapstructure:"redispassword"`
	} `mapstructure:"database"`

	Email struct {
		Address  string `mapstructure:"address"`
		Password string `mapstructure:"password"`
	} `mapstructure:"email"`

	ImageHost struct {
		AccessKey string `mapstructure:"access_key"`
		SecretKey string `mapstructure:"secret_key"`
		Bucket    string `mapstructure:"bucket"`
		Domain    string `mapstructure:"domain"`
	} `mapstructure:"image_host"`
}

var AllConfig Config

// 加载配置
func LoadConfig() {
	viper.SetConfigName("config")   // 文件名
	viper.SetConfigType("yaml")     // 文件类型
	viper.AddConfigPath("./config") // 在 config 目录查找

	// 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("读取配置文件失败: %v", err)
	}

	// 解析到结构体
	if err := viper.Unmarshal(&AllConfig); err != nil {
		log.Fatalf("解析配置文件失败: %v", err)
	}

	fmt.Println("配置加载成功！")
}
