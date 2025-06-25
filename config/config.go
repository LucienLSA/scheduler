package config

import (
	"flag"
	"fmt"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var Conf = new(multipleConfig)

type multipleConfig struct {
	*ServerConfig `mapstructure:"server"`
	*DBConfig     `mapstructure:"db"`
	*LogConfig    `mapstructure:"log"`
}
type ServerConfig struct {
	Addr          string `mapstructure:"addr"`
	Port          int    `mapstructure:"port"`
	Mode          string `mapstructure:"mode"`
	ContextPath   string `mapstructure:"context-path"`
	ConsoleEnable bool   `mapstructure:"console-enable"`
	TLS           bool   `mapstructure:"tls"`
	CertFile      string `mapstructure:"cert-file"`
	KeyFile       string `mapstructure:"key-file"`
}

type LogConfig struct {
	Level      string `mapstructure:"level"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxAge     int    `mapstructure:"max_age"`
	MaxBackups int    `mapstructure:"max_backups"`
	Filename   string `mapstructure:"filename"`
	FilePath   string `mapstructure:"filepath"`
}
type DBConfig struct {
	Dsn              string `mapstructure:"dsn"`
	ExecutedLockTime int64  `mapstructure:"executed-lock-time"`
}

func Init() (err error) {
	var filePath string
	flag.StringVar(&filePath, "filePath", "./config/config.yaml", "配置文件")
	//解析命令行参数
	flag.Parse()
	// fmt.Println(filePath)
	// fmt.Println(flag.Args())
	// fmt.Println(flag.NArg())
	// fmt.Println(flag.NFlag())
	viper.SetConfigFile(filePath)

	err = viper.ReadInConfig()
	if err != nil {
		// 读取配置信息失败
		fmt.Printf("viper.ReadInConfig failed, err:%v\n", err)
		return
	}
	// 配置信息的反序列化
	if err := viper.Unmarshal(Conf); err != nil {
		fmt.Printf("viper Unmarshal failed, err:%v\n", err)
	}
	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		fmt.Println("配置文件被修改了")
		if err := viper.Unmarshal(Conf); err != nil {
			fmt.Printf("viper Unmarshal failed, err:%v\n", err)
		}
	})
	return
}
