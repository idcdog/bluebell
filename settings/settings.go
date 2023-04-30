package settings

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

// 全局变量，保存配置信息
var Conf = new(AppConfig)

type AppConfig struct {
	Name         string `mapstructure:"name"`
	Mode         string `mapstructure:"mode"`
	Version      string `mapstructure:"version"`
	Port         int    `mapstructure:"port"`
	StartTime    string `mapstructure:"start_time"`
	MachineID    int64  `mapstructure:"machine_id"`
	Page         int64  `mapstructure:"page"`
	Size         int64  `mapstructure:"size"`
	MaxPageSize  int64  `mapstructure:"max_page_size"`
	*LogConfig   `mapstructure:"log"`
	*MySQLConfig `mapstructure:"mysql"`
	*RedisConfig `mapstructure:"redis"`
}

type LogConfig struct {
	Level      string `mapstructure:"level"`
	Filename   string `mapstructure:"filename"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxAge     int    `mapstructure:"max_age"`
	MaxBackups int    `mapstructure:"max_backups"`
	Compress   bool   `mapstructure:"compress"`
}
type MySQLConfig struct {
	Host         string `mapstructure:"host"`
	User         string `mapstructure:"user"`
	Password     string `mapstructure:"password"`
	DbName       string `mapstructure:"dbname"`
	Port         int    `mapstructure:"port"`
	MaxOpenConns int    `mapstructure:"max_open_conns"`
	MaxIdleConns int    `mapstructure:"max_idle_conns"`
}
type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
	Port     int    `mapstructure:"port"`
	PoolSize int    `mapstructure:"pool_size"`
}

func Init() (err error) {
	viper.SetConfigName("config")
	// viper.SetConfigName("config") // 用于指定文件名（不包含后缀）， 可能找到的是config.json，也可能是config.yaml
	// viper.SetConfigType("yaml") // 只用于远程读取配置文件后以什么格式进行解析
	viper.AddConfigPath(".")
	viper.AddConfigPath("./conf")
	binPath, _ := filepath.Abs(os.Args[0])
	viper.AddConfigPath(filepath.Base(binPath))
	viper.AddConfigPath(filepath.Join(filepath.Base(binPath), "conf"))
	err = viper.ReadInConfig()
	if err != nil {
		return
	}
	if err = viper.Unmarshal(Conf); err != nil {
		fmt.Printf("viper.Unmarshal failed, err: %v\n", err)
	}
	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		fmt.Println("config file changed!")
		if err = viper.Unmarshal(Conf); err != nil {
			fmt.Printf("viper.Unmarshal failed, err: %v\n", err)
		}
	})

	return
}
