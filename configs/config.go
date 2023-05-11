package configs

import (
	"fmt"
	"github.com/spf13/viper"
)

var Conf = new(Config)

// App 服务基础信息
type App struct {
	Name    string `mapstructure:"name"`
	Port    string `mapstructure:"port"`
	SerPort string `mapstructure:"serPort"`
}
type Config struct {
	App App `mapstructure:"app"` // 服务基本配置
}

func init() {
	viper.SetConfigFile("./configs/dev.conf.yaml")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Sprintf("读取配置文件失败，原因：%s\n", err.Error()))
	}
	if err := viper.Unmarshal(Conf); err != nil {
		panic(fmt.Sprintf("解析配置文件失败，原因：%s\n", err.Error()))
	}
}
