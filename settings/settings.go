package settings

import (
	"log"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

func Init() error {
	viper.SetConfigFile("config.json")
	//viper.SetConfigName("config") // 配置文件名称(无扩展名)
	//viper.SetConfigType("json")   // 如果配置文件的名称中没有扩展名，则需要配置此项
	viper.AddConfigPath(".")    // 还可以在工作目录中查找配置
	err := viper.ReadInConfig() // 查找并读取配置文件
	if err != nil {             // 处理读取配置文件的错误
		return err
	}
	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		log.Println("config changed...")
	})
	return nil
}
