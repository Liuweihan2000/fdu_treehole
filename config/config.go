package config

import (
	"GoProject/fudan_bbs/utils"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"log"
)

func InitConfigFile() {
	viper.SetConfigType("yml")
	viper.AddConfigPath(".")
	viper.SetConfigFile("config.yml")
	err := viper.ReadInConfig()
	utils.FatalErrorHandle(err, "error occured while reading config file")

	// 实时监测配置文件的改动，无需重启服务器
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		log.Println("Config file change detected: ", e.Name)
	})
}


